package grpc

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
)

type ServerRecvMsgInterceptor struct {
}

func (h *ServerRecvMsgInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	fmt.Printf("[Interceptor] serverRecvMsg\n")

	cs := invocation.CallerInstance().(*nativeserverStream)
	ctx := cs.Context()
	md, _ := metadata.FromIncomingContext(ctx)

	peerKey, ok := peer.FromContext(ctx)
	remoteAddr := ""
	if ok {
		remoteAddr = peerKey.Addr.String()
	}
	fmt.Printf("[remoteAddr] :%v\n", remoteAddr)
	if remoteAddr != "127.0.0.1:11800" {

		s, err := tracing.CreateEntrySpan("RecvMsg", func(headerKey string) (string, error) {
			Value := ""
			vals := md.Get(headerKey)
			if len(vals) > 0 {
				Value = vals[0]
			}
			fmt.Printf("[Key/Value] :\n[KEY] :%v\n[Value] :%v\n", headerKey, Value)
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

func (h *ServerRecvMsgInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	if invocation.GetContext() != nil {
		span := invocation.GetContext().(tracing.Span)

		if err, ok := result[0].(error); ok && err != nil {
			span.Error(err.Error())
		}

		span.End()
	}

	return nil
}

/*
skywalking-go 2023/07/24 12:29:39
execute interceptor before invoke error,
instrument name: grpc,
interceptor name: ServerRecvMsgInterceptor,
function ID: google_golang_org_grpc_serverStreamRecvMsg,
error: interface
conversion: interface {} is *grpc.serverStream, not *grpc.clientStream,
*/
