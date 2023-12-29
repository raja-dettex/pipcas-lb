package encoding

import (
	"fmt"
	"io"
	"strings"

	"github.com/raja-dettex/pipcas-lb/messages"
)

type Parser interface {
	Decode(io.Reader, *messages.RPC) error
}

type ConnectionParser struct {
}

func (parser *ConnectionParser) Decode(r io.Reader, rpc *messages.RPC) error {
	for {
		buff := make([]byte, 2048)
		n, err := r.Read(buff)
		if err == io.EOF {
			fmt.Println("Reached end of stream")
			return nil
		}
		if err != nil {
			return nil
		}
		str := strings.Split(string(buff[:n]), " ")
		if str[0] == "read" {

			rpc.Header = messages.RPCMessageRead
		} else if str[0] == "write" {
			rpc.Header = messages.RPCMessageWrite
		} else if str[0] == "add" {
			rpc.Header = messages.RPCMessageAddToPool
		}
		rpc.Payload = str[1:]
		return nil
	}
}
