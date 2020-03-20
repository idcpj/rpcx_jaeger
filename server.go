package main

import (
	"github.com/idcpj/rpcx_jaeger/lib"
	"github.com/smallnest/rpcx/server"
)

func main() {
	closer := lib.InitJaeger("rpc_server")
	defer closer.Close()

	s := server.NewServer()
	//s.Register(new(Arith), "")
	//两种注册方式  与下同
	err := s.Register(new(lib.Args), "")
	if err != nil {
		panic(err)
	}
	err = s.Serve("tcp", "127.0.0.1:8972")
	if err != nil {
		panic(err)
	}
}
