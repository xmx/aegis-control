package telemetry

import (
	"net/http"

	oteltrace "go.opentelemetry.io/otel/trace"
)

type Config struct {
	Endpoint    string            // 接入端点
	Insecure    bool              // 是否跳过证书校验
	Headers     map[string]string // Header
	ServiceName string            // 服务名
	HTTPClient  *http.Client      // 底层 HTTPClient（可选）
}

type TraceProvider interface {
	oteltrace.TracerProvider

	// Reconfigure 动态修改配置。
	Reconfigure(cfg Config) (old oteltrace.TracerProvider, err error)

	// Disable 关闭上报。
	Disable()
}

type traceProvider struct {
	next oteltrace.TracerProvider
	exp  *spanExporter
	proc *spanProcessor
}
