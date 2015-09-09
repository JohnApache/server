package agent

import (
	"bufio"
	"encoding/binary"
	"net"
	"sync"
)

type ConnNet struct {
	net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mutex  sync.Mutex
}

func NewConnNet(conn net.Conn) *ConnNet {
	if conn == nil {
		return nil
	}
	ret := ConnNet{
		Conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
	return &ret
}

func (conn *ConnNet) ReadMsg() (b []byte, err error) {
	head := []byte{0, 0}
	_, err = conn.reader.Read(head)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint16(head)
	body := make([]byte, size)
	sum := uint16(0)
	i := 0
	for {
		i, err = conn.reader.Read(body[sum:])
		if err != nil {
			return nil, err
		}
		if sum += uint16(i); sum == size {
			break
		}
	}
	return body, err
}

func (conn *ConnNet) WriteMsg(b []byte) (err error) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	size := uint16(len(b))
	err = binary.Write(conn.writer, binary.BigEndian, size)
	if err != nil {
		return err
	}
	_, err = conn.writer.Write(b)
	if err != nil {
		return err
	}
	return conn.writer.Flush()
}

func (conn *ConnNet) LocalAddr() string {
	return conn.Conn.LocalAddr().String()
}
func (conn *ConnNet) RemoteAddr() string {
	return conn.Conn.RemoteAddr().String()
}
