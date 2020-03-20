package main

import (
	"context"
	"github.com/idcpj/rpcx_jaeger/lib"
	"github.com/smallnest/rpcx/client"
	"log"
	"time"
)

func main() {
	closer := lib.InitJaeger("rpc2")
	defer closer.Close()

	addrs := []string{"127.0.0.1:8972"}
	var address []*client.KVPair
	for _, addr := range addrs {
		address = append(address, &client.KVPair{Key: addr})
	}

	d := client.NewMultipleServersDiscovery(address)
	xclient := client.NewXClient("Args", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()
	time.Sleep(1 * time.Second)
	rpcx2(xclient)
}

func rpcx2(xclient client.XClient) {
	span, ctx, err := lib.GenerateSpanWithContext(context.Background(), "first2")
	if err != nil {
		panic(err)
	}
	defer span.Finish()

	time.Sleep(1 * time.Second)
	type Args struct {
		N, M int
	}
	type Res struct {
		Cal int
	}

	args := Args{7, 8}
	res := Res{}

	err = xclient.Call(ctx, "Multiply", &args, &res)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	log.Printf("%+v", res)
}
