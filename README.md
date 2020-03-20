## 使用
```

docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest


go get  ./...
go run client.go
go run server.go

```