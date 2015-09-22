package agent

import (
	"errors"
	"fmt"
	"time"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server/cfg"
)

type sessions map[int]map[uint]*Session

var sesss = sessions{}

func (s sessions) Sync(se *Session) *Session {

	if s[se.SerId] == nil {
		s[se.SerId] = map[uint]*Session{}
	}
	uniq := se.ToUint()

	if s[se.SerId][uniq] == nil {
		s[se.SerId][uniq] = se
	} else if se != nil {
		s[se.SerId][uniq].copys(se)
	}

	return s[se.SerId][uniq]
}

func (s sessions) Leave(se *Session) {
	if s[se.SerId] == nil {
		return
	}
	delete(s[se.SerId], se.ToUint())
}

type Session struct {
	base.Unique
	Data           *base.Stora
	Rooms          *base.Stora
	ConnectTime    time.Time
	LastPacketTime time.Time
	Dirtycount     uint
	SerId          int
	tmp            map[string]*base.EncodeBytes
}

func NewSession() *Session {
	s := Session{
		ConnectTime:    time.Now(),
		LastPacketTime: time.Now(),
		Dirtycount:     0,
		SerId:          cfg.SelfId,
		Data:           base.NewStora(),
		Rooms:          base.NewStora(),
	}
	s.InitUint()
	return &s
}
func (s *Session) copys(e *Session) {
	s.Data = e.Data
	s.Rooms = e.Rooms
	s.ConnectTime = e.ConnectTime
	s.LastPacketTime = e.LastPacketTime
	s.Dirtycount = e.Dirtycount
}
func (s *Session) refresh() {
	s.LastPacketTime = time.Now()
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

//func (s *Session) push(reply interface{}) (err error) {
//	return s.Push(reply, nil)
//}

func (s *Session) Push(reply interface{}, hand []byte) (err error) {
	return s.Send(&Response{
		Response: base.EnJson(reply),
		Head:     hand,
	})
}

func (s *Session) Sync() *Session {
	return sesss.Sync(s)

}

func (s *Session) Already(args Request, reply *Response, f func()) {
	f()
	reply.Session = s

}

func (s *Session) Mutex(f func()) {
	var lockreply LockResponse
	err := cfg.GetServer(s.SerId).Client().Call("Connect.Lock", LockRequest{
		Uniq: s.ToUint(),
		Hold: cfg.SelfId,
	}, &lockreply)
	if err != nil {
		return
	}

	s = lockreply.Session.Sync()
	defer func() {
		var unlockreply Response
		unlockreply.Session = s
		cfg.GetServer(s.SerId).Client().Send("Connect.Unlock", UnlockRequest{
			Uniq:  s.ToUint(),
			Reply: &unlockreply,
		})
	}()
	f()
}

func (s *Session) NonSync(f func()) {
	f()
}

func (s *Session) RoomsUniq(name string) uint {
	d := data{}
	s.Rooms.Get(name, &d)
	return d.Id
}

func (s *Session) RoomsHead(name string) []byte {
	d := data{}
	s.Rooms.Get(name, &d)
	return d.Head
}

//func (s *Session) SumData(i interface{}) {
//	s.Data.SumJson(base.EnJson(i))
//}

//func (s *Session) DeData(i interface{}) {
//	s.Data.DeJson(i)
//}

//func (s *Session) EnData(i interface{}) {
//	s.Data.EnJson(i)
//}

//func (s *Session) GetRoomsData() map[string]interface{} {
//	return s.Rooms.Data()
//}

//func (s *Session) Set(k string, v interface{}) {
//	if s.tmp == nil {
//		s.tmp = map[string]*base.EncodeBytes{}
//		s.DeData(&s.tmp)
//		if s.tmp == nil {
//			s.tmp = map[string]*base.EncodeBytes{}
//		}
//	}
//	s.tmp[k] = base.EnJson(v)
//}

//func (s *Session) Get(k string, v interface{}) {
//	if s.tmp == nil {
//		s.tmp = map[string]*base.EncodeBytes{}
//		s.DeData(&s.tmp)
//	}
//	if s.tmp[k] != nil {
//		s.tmp[k].DeJson(v)
//	}
//}
