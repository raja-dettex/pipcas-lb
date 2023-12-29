package messages

import "net"

type RPCMessageHeader byte

const (
	RPCMessageRead      RPCMessageHeader = 0x01
	RPCMessageWrite     RPCMessageHeader = 0x02
	RPCMessageAddToPool RPCMessageHeader = 0x03
)

type RPC struct {
	Addr    net.Addr
	Conn    net.Conn
	Header  RPCMessageHeader
	Payload []string
}
