package lib

import (
	"context"
	"time"
)

type Args struct {
	N, M int
}
type Res struct {
	Cal int
}

func (t *Args) Multiply(ctx context.Context, args *Args, reply *Res) error {
	span, _, err2 := GenerateSpanWithContext(ctx, "rpcx1")
	if err2 != nil {
		panic(err2)
	}
	defer span.Finish()
	reply.Cal = args.N * args.M
	time.Sleep(100 * time.Millisecond)
	return nil
}
