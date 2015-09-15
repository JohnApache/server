package agent

import (
	"errors"
	"fmt"
	"time"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server/cfg"
)

type sessions map[int]map[uint]*Session

func newSessions() sessions {
	return sessions{}
}

func (s sessions) Sync(se *Session) {
	if s[se.SerId] == nil {
		s[se.SerId] = map[uint]*Session{}
	}
	if s[se.SerId][se.ToUint()] == nil || s[se.SerId][se.ToUint()].LastPacketTime.UnixNano() < se.LastPacketTime.UnixNano() {
		s[se.SerId][se.ToUint()] = se
	} else {
		se = s[se.SerId][se.ToUint()]
	}

}

type Session struct {
	base.Unique
	Data           *base.EncodeBytes
	ConnectTime    time.Time
	LastPacketTime time.Time
	Dirtycount     uint
	SerId          int
	occupy         chan func()
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

func (s *Session) Mutex(f func()) {
	if s.occupy == nil {
		s.Refresh()
		var lockreply LockResponse
		err := cfg.GetServer(s.SerId).Client().Call("Connect.Lock", LockRequest{
			Uniq: s.ToUint(),
			Hold: cfg.SelfId,
		}, &lockreply)
		if err != nil {
			return
		}
		*s = *lockreply.Session

		var unlockreply Response
		defer func() {
			unlockreply.Coverage = s.Data
			err = cfg.GetServer(s.SerId).Client().Send("Connect.Unlock", UnlockRequest{
				Uniq:  s.ToUint(),
				Reply: &unlockreply,
			})
		}()
	}
	s.Occupy(f)
}

func (s *Session) Occupy(f func()) {
	if s.occupy == nil {
		s.occupy = make(chan func(), 10)
		s.occupy <- f
		defer func() {
			close(s.occupy)
			s.occupy = nil
		}()
		for {
			select {
			case v, ok := <-s.occupy:
				if ok {
					v()
				} else {
					return
				}
			default:
				return
			}
		}
	} else {
		s.occupy <- f
	}
	return
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
		Hold: cfg.SelfId,
	}, &reply)
	if err != nil {
		return err
	}
	*s = *reply.Session
	return nil
}

func (s *Session) MutexSession(f func() *Response) (err error) {

	//s.occupy = true
	defer func() {
		//s.occupy = false
		if x := recover(); x != nil {
			err = errors.New("Session.MutexSession: " + fmt.Sprintln(x))
		}
	}()

	var reply LockResponse
	err = cfg.GetServer(s.SerId).Client().Call("Connect.Lock", LockRequest{
		Uniq: s.ToUint(),
		Hold: cfg.SelfId,
	}, &reply)
	if err != nil {
		return err
	}
	*s = *reply.Session
	return cfg.GetServer(s.SerId).Client().Send("Connect.Unlock", UnlockRequest{
		Uniq:  s.ToUint(),
		Reply: f(),
	})
}
func (s *Session) Sum(i interface{}) {
	//if s.occupy {
	s.Data = base.SumJson(s.Data, base.EnJson(i))
	//}

}

//func (s *Session) LockSession() (err error) {
//	defer func() {
//		if x := recover(); x != nil {
//			err = errors.New("Session.LockSession: " + fmt.Sprintln(x))
//		}
//	}()
//	var reply LockResponse
//	err = cfg.GetServer(s.SerId).Client().Call("Connect.Lock", LockRequest{
//		Uniq: s.ToUint(),
//		Hold: cfg.SelfId,
//	}, &reply)
//	if err != nil {
//		return err
//	}
//	*s = *reply.Session
//	return nil
//}

//func (s *Session) UnlockSession(reply *Response) error {
//	return cfg.GetServer(s.SerId).Client().Send("Connect.Unlock", UnlockRequest{
//		Uniq:  s.ToUint(),
//		Reply: reply,
//	})
//}

func (s *Session) Change(i interface{}) error {
	return cfg.GetServer(s.SerId).Client().Send("Connect.Change", ChangeRequest{
		Uniq: s.ToUint(),
		Data: base.EnJson(i),
	})
}
