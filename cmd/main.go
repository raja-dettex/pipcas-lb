package main

import (
	"log"
	"os"

	"github.com/raja-dettex/pipcas-lb/encoding"
	"github.com/raja-dettex/pipcas-lb/messages"
	"github.com/raja-dettex/pipcas-lb/server"
)

var (
	listenAddr = os.Getenv("LISTEN_ADDR")
)

func main() {
	pool := &server.ServerPool{HashMap: make(map[uint64]*server.Server)}
	rpcCh := make(chan *messages.RPC)
	opts := server.LBOpts{ListenAddr: listenAddr, Pool: pool, Parser: &encoding.ConnectionParser{}, RPCCh: rpcCh}
	lbServer := server.NewLBServer(opts)
	if err := lbServer.Start(); err != nil {
		log.Fatal(err)
	}
	select {}
}
