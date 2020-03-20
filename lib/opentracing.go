package lib

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"log"
)

/**
改函数只适用于  jaeger
通过传递 tracekey string 来追踪
@tracekey 格式 : 6f8b8a1101b0124f:6f8b8a1101b0124f:0000000000000000:1

首个 span
span, carrier, _ := lib.GenerateSpanWithContext( "first2","")
defer span.Finish()

之后 通过传递 carrier 来实现追踪
span, _, _ := lib.GenerateSpanWithContext( "first2",carrier)
defer span.Finish()
*/
func GenerateSpanWithContext(operationName string, traceKey string) (opentracing.Span, string, error) {
	var span opentracing.Span

	tracer := opentracing.GlobalTracer()

	if traceKey != "" {
		spanContext, err := jaeger.ContextFromString(traceKey) //通过 traceKey 获取 spanContext
		if err != nil {
			log.Printf("metadata error %s\n", err)
			return nil, "", err
		}
		span = tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
	} else {
		span = opentracing.StartSpan(operationName)
	}

	metadata := opentracing.TextMapCarrier(make(map[string]string))
	err := tracer.Inject(span.Context(), opentracing.TextMap, metadata)
	if err != nil {
		return nil, "", err
	}
	return span, metadata[jaeger.TraceContextHeaderName], nil
}

func InitJaeger(service string) io.Closer {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const", //全部采样
			Param: 1,       //1 开启全部采样,0 关闭全部采样,可通过 环境变量 JAEGER_SAMPLER_PARAM 控制
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.0.110:6831", //web查看 ip:16686
		},
	}
	tracer, closer, err := cfg.NewTracer() //log.StdLogger 只要实现日志接口即可
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger: %v\n", err)
	}
	opentracing.SetGlobalTracer(tracer)

	return closer
}
