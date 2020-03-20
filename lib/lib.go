package lib

import (
	"context"
	"time"
)

type Args struct {
	N, M    int
	Carrier string
}
type Res struct {
	Cal     int
	Carrier string
}

func (t *Args) Multiply(ctx context.Context, args *Args, reply *Res) error {
	span, carrier, err2 := GenerateSpanWithContext("rpcx1", args.Carrier)
	if err2 != nil {
		panic(err2)
	}
	defer span.Finish()

	reply.Cal = args.N * args.M
	reply.Carrier = carrier
	time.Sleep(100 * time.Millisecond)
	return nil
}
