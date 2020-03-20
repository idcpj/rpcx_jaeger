package lib

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/smallnest/rpcx/share"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"log"
)

// 只适用于 jaeger
// 原理用 context 传递 "__req_metadata":""uber-trace-id -> 6f8b8a1101b0124f:6f8b8a1101b0124f:0000000000000000:1
// 										uber-trace-id traceID : spanID : parentID: sampled bool
func GenerateSpanWithContext(ctx context.Context, operationName string) (opentracing.Span, context.Context, error) {
	md := ctx.Value(share.ReqMetaDataKey) // share.ReqMetaDataKey 固定值 "__req_metadata"  可自定义
	var span opentracing.Span

	tracer := opentracing.GlobalTracer()

	if md != nil {
		carrier := opentracing.TextMapCarrier(md.(map[string]string))
		spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Printf("metadata error %s\n", err)
			return nil, nil, err
		}
		span = tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
	} else {
		span = opentracing.StartSpan(operationName)
	}

	metadata := opentracing.TextMapCarrier(make(map[string]string))
	err := tracer.Inject(span.Context(), opentracing.TextMap, metadata)
	if err != nil {
		return nil, nil, err
	}
	//把metdata 携带的 traceid,spanid,parentSpanid 放入 context
	ctx = context.WithValue(context.Background(), share.ReqMetaDataKey, (map[string]string)(metadata))
	return span, ctx, nil
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
