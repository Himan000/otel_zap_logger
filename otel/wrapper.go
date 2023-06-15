package otel

/*
暂时只用到otel trace的功能，未用到zap的log功能。
*/

import (
	"context"
	"net/http"

	"github.com/Himan000/otel_zap_logger"
	"github.com/Himan000/otel_zap_logger/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var TRACE_ID = "trace_id"

type Logger struct {
	ctx context.Context
}

var Default = New()

func New() *Logger {
	return &Logger{}
}

/*
	初始时otel，包括Tempo服务器，以及在trace显示的信息。
*/
func (l *Logger) Init(g *gin.Engine) *Logger {
	//项目配置
	c := config.New(viper.GetViper())
	_ = c.Load()

	conf := otel_zap_logger.Config{
		Debug:       true,
		EnableTrace: true,
		// EnableLog:   true,
		// File:               "./" + viper.GetString("APP_ID") + ".log",
		// File:               "./run.log",
		TracerProviderType: "jaeger",
		TraceSampleRatio:   1,
		JaegerServer:       viper.GetString(config.JAEGER_SERVER), // 1
	}

	// 这些信息只会在trace显示
	otel_zap_logger.Init(conf,
		otel_zap_logger.String("service.name", viper.GetString(config.APP_ID)),
		otel_zap_logger.String("envtype", viper.GetString(config.ENV_TYPE)),
		otel_zap_logger.String("appid", viper.GetString(config.APP_ID)),
	)

	g.Use(GinMiddleware(viper.GetString(config.APP_ID))) // 用于拼凑请求全路径

	return l
}

func (l *Logger) Start(ctx context.Context) context.Context {
	l.ctx = otel_zap_logger.Start(ctx, viper.GetString(config.APP_ID))
	return l.ctx
}

func Start(c context.Context) (context.Context, string) {
	ctx := Default.Start(c)
	traceId := GetTraceId(ctx)
	return ctx, traceId
}

func End() {
	Default.end()
}

func (l *Logger) end() {
	otel_zap_logger.End(l.ctx)
}

func GinMiddleware(service string) gin.HandlerFunc {
	return otel_zap_logger.GinMiddleware(service)
}

func HttpInject(ctx context.Context, request *http.Request) {
	otel_zap_logger.HttpInject(ctx, request)
}

func String(key string, value string) otel_zap_logger.Field {
	return otel_zap_logger.String(key, value)
}

func GetTraceId(ctx context.Context) string {
	return otel_zap_logger.TraceID(ctx)
}
