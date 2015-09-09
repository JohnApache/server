package agent

import (
	"errors"
	"fmt"
	"time"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server/cfg"
)

type Session struct {
	base.Unique
	Data           *base.EncodeBytes
	ConnectTime    time.Time
	LastPacketTime time.Time
	Dirtycount     uint
	SerId          int
}

func NewSession() *Session {
	s := Session{
		ConnectTime:    time.Now(),
		LastPacketTime: time.Now(),
		Dirtycount:     0,
		SerId:          cfg.SelfId,
	}
	s.InitUint()
	s.Data = base.EnJson(map[string]uint{
		"none": 0,
	})
	return &s
}

func (s *Session) Refresh() {
	s.LastPacketTime = time.Now()
}

func (s *Session) Push(reply interface{}) (err error) {
	return s.Send(&Response{
		Response: base.EnJson(reply),
	})
}

func (s *Session) Send(reply *Response) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("Session.Send: " + fmt.Sprintln(x))
		}
	}()
	return cfg.GetServer(s.SerId).Client().Send("Connect.Push", PushRequest{
		Uniq:  s.ToUint(),
		Reply: reply,
	})
}

func (s *Session) SyncSession() (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("Session.SyncSession: " + fmt.Sprintln(x))
		}
	}()
	var reply LockResponse
	err = cfg.GetServer(s.SerId).Client().Call("Connect.Sync", LockRequest{
		Uniq: s.ToUint(),
	}, &reply)
	if err != nil {
		return err
	}
	*s = *reply.Session
	return nil
}

func (s *Session) LockSession() (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("Session.LockSession: " + fmt.Sprintln(x))
		}
	}()
	var reply LockResponse
	err = cfg.GetServer(s.SerId).Client().Call("Connect.Lock", LockRequest{
		Uniq: s.ToUint(),
	}, &reply)
	if err != nil {
		return err
	}
	*s = *reply.Session
	return nil
}

func (s *Session) UnlockSession(reply *Response) error {
	return cfg.GetServer(s.SerId).Client().Send("Connect.Unlock", UnlockRequest{
		Uniq:  s.ToUint(),
		Reply: reply,
	})
}

func (s *Session) Change(i interface{}) error {
	return cfg.GetServer(s.SerId).Client().Send("Connect.Change", ChangeRequest{
		Uniq: s.ToUint(),
		Data: base.EnJson(i),
	})
}
