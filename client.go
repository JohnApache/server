package server

import (
	"net/rpc"

	"github.com/wzshiming/base"
)

type Client struct {
	rpc.Client
}

var buf = make(map[string]*Client)

func NewClient(addr string) *Client {
	if conn := buf[addr]; conn != nil {
		return conn
	}
	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		base.ERR(err)
		return nil
	}
	conn := &Client{
		Client: *c,
	}
	buf[addr] = conn
	return conn
}

func (cl *Client) Send(serviceMethod string, args interface{}) error {
	var t int
	return cl.Call(serviceMethod, args, &t)
}

func (cl *Client) Take(serviceMethod string, reply interface{}) error {
	return cl.Call(serviceMethod, 0, reply)
}

func (cl *Client) TakeClasss() (c Classs, err error) {
	err = cl.Take("Classs.Take", &c)
	return
}

func (cl *Client) ShutdownNow() error {
	return cl.Send("Shutdown.Now", 0)
}
