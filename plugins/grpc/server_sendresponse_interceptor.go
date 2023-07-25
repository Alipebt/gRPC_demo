package grpc

import (
	"fmt"
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type ServerSendResponseInterceptor struct {
}

func (h *ServerSendResponseInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	fmt.Printf("[Interceptor] serverSendResponse\n")

	cs := invocation.Args()[1].(*nativeStream)

	ctx := cs.Context()
	peerKey, ok := peer.FromContext(ctx)
	remoteAddr := ""
	if ok {
		remoteAddr = peerKey.Addr.String()
	}
	// fmt.Printf("[context] :%v\n", ctx)
	// fmt.Printf("[remoteAddr] :%v\n", remoteAddr)

	if remoteAddr != "127.0.0.1:11800" {

		s, err := tracing.CreateExitSpan("SendResponse", remoteAddr, func(headerKey, headerValue string) error {

			ctx = metadata.AppendToOutgoingContext(ctx, headerKey, headerValue)

			return nil
		},
			tracing.WithLayer(tracing.SpanLayerRPCFramework),
			tracing.WithTag(tracing.TagURL, "SendResponse"),
			tracing.WithComponent(5017),
		)

		if err != nil {
			return err
		}

		invocation.SetContext(s)
	}
	return nil
}

func (h *ServerSendResponseInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	if invocation.GetContext() != nil {

		span := invocation.GetContext().(tracing.Span)

		if err, ok := result[0].(error); ok && err != nil {
			span.Error(err.Error())
		}

		span.End()
	}

	return nil
}
