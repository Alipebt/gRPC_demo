package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strings"

	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
	"google.golang.org/grpc"
)

type ClientStreamingInterceptor struct {
}

func (h *ClientStreamingInterceptor) BeforeInvoke(invocation operator.Invocation) error {

	// 获取参数
	ctx := invocation.Args()[0].(context.Context)
	method := invocation.Args()[2].(string)
	clientconn := invocation.CallerInstance().(*grpc.ClientConn)
	remoteAddr := clientconn.Target()

	if !strings.HasPrefix(method, "/skywalking") {
		fmt.Printf("[Interceptor] NewStream\n")

		s, err := tracing.CreateExitSpan(method, remoteAddr, func(headerKey, headerValue string) error {
			ctx = metadata.AppendToOutgoingContext(ctx, headerKey, headerValue)
			invocation.ChangeArg(0, ctx)

			return nil
		},
			tracing.WithLayer(tracing.SpanLayerRPCFramework),
			tracing.WithTag(tracing.TagURL, method),
			tracing.WithTag("Communication Mode", "Streaming"),
			tracing.WithComponent(5017),
		)

		if err != nil {
			return err
		}

		invocation.SetContext(s)
	}
	return nil
}

func (h *ClientStreamingInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	if invocation.GetContext() != nil {

		span := invocation.GetContext().(tracing.Span)

		if err, ok := result[0].(error); ok && err != nil {
			span.Error(err.Error())
		}

		span.End()
	}
	return nil
}
