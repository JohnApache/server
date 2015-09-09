package agent

import (
	"net"

	"github.com/wzshiming/base"
)

func NewConn(addr string) Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		base.ERR(err)
		return nil
	}
	return NewConnNet(conn)
}
