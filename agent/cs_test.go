package agent

import (
	"net"
	"testing"

	"github.com/wzshiming/base"
)

func Test_CS(t *testing.T) {
	go tserver(t)
	tclient(t)
}

func Test_tcp(t *testing.T) {
	var err error
	msg1 := []byte("hello")
	msg2 := []byte("bye")
	go func() {
		listen, err := net.Listen("tcp", ":7790")
		checkError(err, t)
		conn, err := listen.Accept()
		checkError(err, t)
		stcp := NewConnNet(conn)
		b, err := stcp.ReadMsg()
		checkError(err, t)
		if string(b) != string(msg1) {
			t.Fatal(b, msg1)
		}
		stcp.WriteMsg(msg2)
		defer conn.Close()
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:7790")
	checkError(err, t)
	ctcp := NewConnNet(conn)
	ctcp.WriteMsg(msg1)
	b, err := ctcp.ReadMsg()
	checkError(err, t)
	if string(b) != string(msg2) {
		t.Fatal(b, msg2)
	}
}

func tserver(t *testing.T) {
	ag := NewAgent(16, func(user *User, msg []byte) error {
		base.NOTICE(string(msg))
		user.WriteMsg(msg)
		user.Close()
		return nil
	})
	err := NewListener(7710, func(conn Conn) {
		ag.Join(conn)
	}).Start()
	checkError(err, t)
}

func tclient(t *testing.T) {
	conn := NewConn("127.0.0.1:7710")
	if conn == nil {
		t.Fail()
		return
	}
	bi := []byte("hello")
	conn.WriteMsg(bi)
	b, err := conn.ReadMsg()
	checkError(err, t)
	if string(b) == "" && string(b) != string(bi) {
		t.Fail()
	}
	b, err = conn.ReadMsg()
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
