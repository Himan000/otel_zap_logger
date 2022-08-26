package otel_zap_logger

import (
	"context"
	"net/http"

	"github.com/Himan000/otel_zap_logger/propagation/extract"
	"github.com/Himan000/otel_zap_logger/propagation/inject"
	"github.com/gin-gonic/gin"
)

// HTTPInject inject spanContext
func HttpInject(ctx context.Context, request *http.Request) error {
	return inject.HttpInject(ctx, request)
}

// GinMiddleware extract spanContext
func GinMiddleware(service string) gin.HandlerFunc {
	return extract.GinMiddleware(service)
}
