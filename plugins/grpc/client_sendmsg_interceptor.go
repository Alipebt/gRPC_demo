package grpc

import (
	"fmt"
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type ClientSendMsgInterceptor struct {
}

func (h *ClientSendMsgInterceptor) BeforeInvoke(invocation operator.Invocation) error {

	cs := invocation.CallerInstance().(*nativeclientStream)
	ctx := cs.Context()

	peerKey, ok := peer.FromContext(ctx)
	remoteAddr := ""
	if ok {
		remoteAddr = peerKey.Addr.String()
	}

	if remoteAddr != "127.0.0.1:11800" {

		fmt.Printf("[Interceptor] clientSendMsg\n")

		s, err := tracing.CreateExitSpan("SendMsgasdasdas", remoteAddr, func(headerKey, headerValue string) error {
			ctx = metadata.AppendToOutgoingContext(ctx, headerKey, headerValue)
			cs.ctx = ctx
			fmt.Printf("[Context+key/value] :%v\n[KEY] :%v\n[Value] :%v\n", cs.ctx, headerKey, headerValue)

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

func (h *ClientSendMsgInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	if invocation.GetContext() != nil {

		span := invocation.GetContext().(tracing.Span)

		if err, ok := result[0].(error); ok && err != nil {
			span.Error(err.Error())
		}

		span.End()
	}

	return nil
}
