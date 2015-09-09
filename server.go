package server

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/wzshiming/base"
)

type Server struct {
	server rpc.Server
	listen net.Listener
	port   int
	classs Classs
}

func NewServer(port int) *Server {
	se := Server{
		server: *rpc.NewServer(),
		port:   port,
	}
	if err := se.server.Register(&se.classs); err != nil {
		base.ERR(err)
		return nil
	}

	if err := se.server.Register(NewShutdown(se.Stop)); err != nil {
		base.ERR(err)
		return nil
	}
	return &se
}

func (se *Server) Classs() Classs {
	return se.classs
}

func (se *Server) Start() (err error) {
	base.NOTICE("Server start from port", se.port)
	se.listen, err = net.Listen("tcp", fmt.Sprintf(":%d", se.port))
	if err != nil {
		base.ERR(err)
		return
	}

	for {
		if conn, err := se.listen.Accept(); err == nil {
			go se.server.ServeConn(conn)
			//go h.Server.ServeCodec(jsonrpc.NewServerCodec(conn))
		} else {
			base.ERR(err)
			return err
		}
	}
	base.NOTICE("Server stop")
	return nil
}

func (se *Server) Register(rcvr interface{}) error {
	err := se.server.Register(rcvr)
	if err != nil {
		return err
	}
	m := NewMethods(rcvr)
	if len(m.Methods) != 0 {
		se.classs = append(se.classs, *m)
	}
	return nil
}

func (se *Server) RegisterName(name string, rcvr interface{}) error {
	err := se.server.RegisterName(name, rcvr)
	if err != nil {
		return err
	}
	m := NewMethods(rcvr)
	m.Name = name
	if len(m.Methods) != 0 {
		se.classs = append(se.classs, *m)
	}
	return nil
}

func (se *Server) Stop() {
	se.listen.Close()
}
