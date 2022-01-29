package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"time"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: "127.0.0.1:6831",
			LogSpans:           true,
		},
		ServiceName: "weihang",
	}
	// 输出到屏幕
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	parentSpan := tracer.StartSpan("main")
	span := tracer.StartSpan("go-grpc-web", opentracing.ChildOf(parentSpan.Context()))
	// 模拟调用
	time.Sleep(time.Millisecond * 500)
	defer span.Finish()

	span2 := tracer.StartSpan("go-grpc-goods", opentracing.ChildOf(parentSpan.Context()))
	// 模拟调用
	time.Sleep(time.Millisecond * 1000)
	span2.Finish()

	parentSpan.Finish()
}
