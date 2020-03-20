## 使用
```

docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest


go get  ./...
go run client.go
go run server.go

```

## master 分支
使用 context 进行上下文追踪

## tracing_string 分支
使用 手动传递 string 的方式进行追踪

## 使用 zipkin 
请跳转到 `https://github.com/idcpj/LearnRpcx`