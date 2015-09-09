package agent

import (
	"time"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg([]byte) error
	Close() error
	LocalAddr() string
	RemoteAddr() string
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}
