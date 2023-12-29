package server

import (
	"fmt"
	"io"
	"net"

	"github.com/raja-dettex/pipcas-lb/encoding"
	"github.com/raja-dettex/pipcas-lb/messages"
)

type LoadBalancer struct {
	Pool       *ServerPool
	ListenAddr string
	Listener   net.Listener
	parser     encoding.Parser
	rpcCh      chan *messages.RPC
}

func NewLoadBalancer(pool *ServerPool, listenAddr string, rpcCh chan *messages.RPC, parser encoding.Parser) *LoadBalancer {
	return &LoadBalancer{Pool: pool, ListenAddr: listenAddr, rpcCh: rpcCh, parser: parser}
}

func (lb *LoadBalancer) Listen() error {
	fmt.Println("servr listening to port ", lb.ListenAddr)
	ln, err := net.Listen("tcp", lb.ListenAddr)
	if err != nil {
		return err
	}
	lb.Listener = ln
	return nil
}

func (lb *LoadBalancer) AccepConn() {
	for {
		conn, err := lb.Listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go lb.HandleConn(conn)
		}
	}
}

func (lb *LoadBalancer) HandleConn(conn net.Conn) {
	rpc := &messages.RPC{Addr: conn.RemoteAddr(), Conn: conn}
	if err := lb.parser.Decode(conn, rpc); err != nil {
		fmt.Println(err)
		return
	}
	lb.rpcCh <- rpc
	return
}

func (lb *LoadBalancer) Consume() {
	for {
		select {
		case rpc := <-lb.rpcCh:
			fmt.Println("rpc ", rpc)
			go lb.handleRPC(rpc)
		default:
			continue
		}
	}
}

func (lb *LoadBalancer) handleRPC(rpc *messages.RPC) {

	switch (*rpc).Header {
	case messages.RPCMessageRead:
		res, err := lb.handleReadRPC(rpc)
		defer rpc.Conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		rpc.Conn.Write([]byte(res))
		return
		//fmt.Println(s.Address())
	case messages.RPCMessageWrite:
		loc, err := lb.handleWriteRPC(rpc)
		defer rpc.Conn.Close()
		if err != nil {
			fmt.Println("error ", err)
			return
		}
		rpc.Conn.Write([]byte(fmt.Sprintf("%s", loc)))
		return
		//fmt.Println(s.Address())
		//fmt.Println((*rpc).Payload[1:])
	case messages.RPCMessageAddToPool:
		lb.Pool.AddToPool((*rpc).Payload[0], (*rpc).Payload[1])
		lb.Pool.ListPool()
	}
}

func (lb *LoadBalancer) handleWriteRPC(rpc *messages.RPC) (string, error) {
	var lbToClientBuff = make([]byte, 1024)
	s := lb.Pool.FindAvailableServer((*rpc).Payload[0])
	conn, err := net.Dial("tcp", s.Address())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	content := ""
	for i := 1; i < len(rpc.Payload); i++ {
		content += rpc.Payload[i]
		content += " "
	}
	_, err = conn.Write([]byte(fmt.Sprintf("write %s %s", rpc.Payload[0], content)))
	if err != nil {
		return "", err
	}
	for {
		n, err := conn.Read(lbToClientBuff)
		if err == io.EOF {
			// do nothing
			fmt.Println(err)
			return s.Address(), nil
		}
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("response %s", string(lbToClientBuff[:n]))
	}
}
func (lb *LoadBalancer) handleReadRPC(rpc *messages.RPC) (string, error) {
	var lbToClientBuff = make([]byte, 1024)
	var res = ""
	s := (*rpc).Payload[1]
	fileName := (*rpc).Payload[0]
	conn, err := net.Dial("tcp", s)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(fmt.Sprintf("read %s", fileName)))
	if err != nil {
		return "", err
	}
	for {
		n, err := conn.Read(lbToClientBuff)
		if err == io.EOF {
			// do nothing
			return res, nil
		}
		if err != nil {
			fmt.Println(err)
		}
		res += string(lbToClientBuff[:n])
	}
}
