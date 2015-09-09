package server

import (
	"testing"
)

func Test_cs(t *testing.T) {
	go server(t)
	client(t)
}

type Te struct{}

type Request int
type Response int

func (t *Te) Echo(a int, r *int) error {
	*r = a
	return nil
}

func (te *Te) TT(a Request, r *Response) error {
	return nil
}

func server(t *testing.T) {
	s := NewServer(7800)
	s.Register(&Te{})
	s.RegisterName("Hello", &Te{})
	s.Start()
}

func client(t *testing.T) {
	c := NewClient("127.0.0.1:7800")
	var a int
	c.Call("Te.Echo", 1000, &a)
	if a != 1000 {
		t.Fail()
	}

	cl, err := c.TakeClasss()
	if err != nil || len(cl) == 0 {
		t.Fail()
	}
	t.Log(cl)
}
