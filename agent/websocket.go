package agent

import (
	"encoding/base64"
	"time"

	"github.com/gorilla/websocket"
)

type ConnWeb struct {
	*websocket.Conn
	chann chan []byte
}

func NewConnWeb(conn *websocket.Conn) ConnWeb {
	return ConnWeb{
		Conn:  conn,
		chann: make(chan []byte, 12),
	}
}

func (conn ConnWeb) ReadMsg() ([]byte, error) {
	_, b, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(string(b))
}

func (conn ConnWeb) WriteMsg(b []byte) error {
	conn.chann <- []byte(base64.StdEncoding.EncodeToString(b))
	for len(conn.chann) > 0 {
		if v, ok := <-conn.chann; ok {
			conn.WriteMessage(websocket.TextMessage, v)
		} else {
			return nil
		}
	}
	return nil
}

func (conn ConnWeb) SetDeadline(t time.Time) error {
	return conn.UnderlyingConn().SetDeadline(t)
}

func (conn ConnWeb) LocalAddr() string {
	return conn.Conn.LocalAddr().String()
}

func (conn ConnWeb) RemoteAddr() string {
	return conn.Conn.RemoteAddr().String()
}
