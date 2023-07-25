package grpc

import (
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
)

type ServerStreamingInterceptor struct {
}

func (h *ServerStreamingInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	span := tracing.ActiveSpan()
	// 添加标签
	span.Tag("Communication Mode", "Streaming") // 对应拦截器中选择对应的模式

	return nil
}
func (h *ServerStreamingInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {

	return nil
}
