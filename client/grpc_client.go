package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	hello "jaeger_test/proto"
)

const (
	// gRPC 服务地址
	Address = "0.0.0.0:8700"
)

func main() {
	//生成 tracer
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

	//设置为全局
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	// opt加上全局tracer
	dial, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		panic(err)
	}
	//初始化客户端
	c := hello.NewHelloClient(dial)
	sayHello, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "weihang"})
	if err != nil {
		panic(err)
	}

	fmt.Println(sayHello)
}
