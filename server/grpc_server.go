package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	hello "jaeger_test/proto"
	"net"
)

type HelloSrv struct{}

func (h *HelloSrv) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	fmt.Printf("handle SayHello 接受到name ：%s", in.Name)
	return &hello.HelloResponse{Message: fmt.Sprintf("%s", in.Name)}, nil
}

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
	address := "0.0.0.0:8700"
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	//注册服务
	hello.RegisterHelloServer(s, &HelloSrv{})

	fmt.Printf("start to listen %s", address)

	if err = s.Serve(listen); err != nil {
		panic(err)
	}
}
