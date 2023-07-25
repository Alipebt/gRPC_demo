package grpc

import (
	"fmt"
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type ServerSendMsgInterceptor struct {
}

func (h *ServerSendMsgInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	fmt.Printf("[Interceptor] serverSendMsg\n")

	cs := invocation.CallerInstance().(*nativeserverStream)
	ctx := cs.Context()

	peerKey, ok := peer.FromContext(ctx)
	remoteAddr := ""
	if ok {
		remoteAddr = peerKey.Addr.String()
	}

	fmt.Printf("[remoteAddr] :%v\n", remoteAddr)
	if remoteAddr != "127.0.0.1:11800" {

		s, err := tracing.CreateExitSpan("SendMsg", remoteAddr, func(headerKey, headerValue string) error {
			ctx = metadata.AppendToOutgoingContext(ctx, headerKey, headerValue)

			fmt.Printf("[Context+key/value] :%v\n[KEY] :%v\n[Value] :%v\n", ctx, headerKey, headerValue)

			return nil
		},
			tracing.WithLayer(tracing.SpanLayerRPCFramework),
			tracing.WithTag(tracing.TagURL, "SendMsg"),
			tracing.WithComponent(5017),
		)

		if err != nil {
			return err
		}

		invocation.SetContext(s)
	}
	return nil
}

func (h *ServerSendMsgInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	if invocation.GetContext() != nil {

		cs := invocation.CallerInstance().(*nativeserverStream)
		ctx := cs.Context()
		fmt.Printf("AFTER [Context] :%v\n", ctx)

		span := invocation.GetContext().(tracing.Span)

		if err, ok := result[0].(error); ok && err != nil {
			span.Error(err.Error())
		}

		span.End()
	}

	return nil
}
