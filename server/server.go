package server

import (
	"github.com/raja-dettex/pipcas-lb/encoding"
	"github.com/raja-dettex/pipcas-lb/messages"
)

type LBOpts struct {
	ListenAddr string
	Pool       *ServerPool
	Parser     encoding.Parser
	RPCCh      chan *messages.RPC
}

type LBServer struct {
	opts LBOpts
	lb   *LoadBalancer
}

func NewLBServer(opts LBOpts) *LBServer {
	return &LBServer{opts: opts}
}

func (server *LBServer) Start() error {

	lb := NewLoadBalancer(server.opts.Pool, server.opts.ListenAddr, server.opts.RPCCh, server.opts.Parser)
	server.lb = lb
	go lb.Consume()
	if err := server.lb.Listen(); err != nil {
		return err
	}
	go lb.AccepConn()
	return nil
}
