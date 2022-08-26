package otel_zap_logger

import (
	"context"
	"net/http"

	"gitee.com/Himan000/otel_zap_logger/propagation/extract"
	"gitee.com/Himan000/otel_zap_logger/propagation/inject"
	"github.comHiman000
)

// HTTPInject inject spanContext
func HttpInject(ctx context.Context, request *http.Request) error {
	return inject.HttpInject(ctx, request)
}

// GinMiddleware extract spanContext
func GinMiddleware(service string) gin.HandlerFunc {
	return extract.GinMiddleware(service)
}
