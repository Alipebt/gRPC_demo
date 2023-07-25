package grpc

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
)

type ClientRecvMsgInterceptor struct {
}

func (h *ClientRecvMsgInterceptor) BeforeInvoke(invocation operator.Invocation) error {

	cs := invocation.CallerInstance().(*nativeclientStream)
	ctx := cs.Context()

	md, _ := metadata.FromIncomingContext(ctx)

	peerKey, ok := peer.FromContext(ctx)
	remoteAddr := ""
	if ok {
		remoteAddr = peerKey.Addr.String()
	}

	if remoteAddr != "127.0.0.1:11800" {
		fmt.Printf("[Interceptor] clientRecvMsg\n")

		s, err := tracing.CreateEntrySpan("RecvMsg", func(headerKey string) (string, error) {
			Value := ""
			vals := md.Get(headerKey)
			if len(vals) > 0 {
				Value = vals[0]
			}
			// fmt.Printf("[Key/Value] :\n[KEY] :%v\n[Value] :%v\n", headerKey, Value)

			return Value, nil
		},
			tracing.WithLayer(tracing.SpanLayerRPCFramework),
			tracing.WithTag(tracing.TagURL, "RecvMsg"),
			tracing.WithComponent(5017),
		)

		if err != nil {
			return err
		}

		invocation.SetContext(s)
	}
	return nil
}

func (h *ClientRecvMsgInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	if invocation.GetContext() != nil {
		span := invocation.GetContext().(tracing.Span)

		if err, ok := result[0].(error); ok && err != nil {
			span.Error(err.Error())
		}

		span.End()
	}

	return nil
}
