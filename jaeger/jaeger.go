/**
 * @Author:huangshan
 * @Description:
 * @Date: 2022/10/26 10:24
 */
package jaeger

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerzap "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
	"io"
	"time"
)

func InitWithZap(service, agentHost string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &jaegercfg.Configuration{
		ServiceName: service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, //常量
			Param: 1,                       //全部采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			QueueSize: 1024,
			LogSpans:  true,
			//CollectorEndpoint: "http://host.docker.internal:14268/api/traces",
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
			//LocalAgentHostPort:  agentHost,//这个不通
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaegerzap.NewLogger(logger)))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

// 开启一个span 并将traceID spanid设置tag
func StartSpan(ctx context.Context, opName string) (span opentracing.Span, childCtx context.Context) {
	span, childCtx = opentracing.StartSpanFromContext(ctx, opName)
	traceid, spanid := GetTraceIdAndSpanId(span)
	span.SetTag("traceid", traceid)
	span.SetTag("spanid", spanid)
	return
}

func GetTraceIdAndSpanId(span opentracing.Span) (traceid, spanid string) {
	if span == nil {
		return "", ""
	}
	if spanCtx, ok := span.Context().(jaeger.SpanContext); ok {
		traceid = spanCtx.TraceID().String()
		spanid = spanCtx.SpanID().String()
	}
	return
}
