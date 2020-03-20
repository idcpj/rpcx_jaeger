package main

import (
	"context"
	"fmt"
	"github.com/idcpj/rpcx_jaeger/lib"
	"github.com/smallnest/rpcx/client"
	"log"
	"time"
)

func main() {
	closer := lib.InitJaeger("rpc3")
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
	span, carrier, err := lib.GenerateSpanWithContext("first2", "")
	if err != nil {
		panic(err)
	}
	defer span.Finish()

	time.Sleep(1 * time.Second)

	args := lib.Args{7, 8, carrier}
	res := lib.Res{}

	err = xclient.Call(context.Background(), "Multiply", &args, &res)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	fuc3(res)
}

func fuc3(res lib.Res) {
	span, _, err := lib.GenerateSpanWithContext("first2", res.Carrier)
	if err != nil {
		panic(err)
	}
	defer span.Finish()

	time.Sleep(1 * time.Second)
	fmt.Printf("%+v\n", res.Cal)
}
