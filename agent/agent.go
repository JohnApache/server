package agent

import (
	"sync"
	"time"

	"github.com/wzshiming/base"
)

type User struct {
	Conn
	sync.RWMutex
	Session *Session
	outtime time.Duration
}

func NewUser(sess *Session, conn Conn) *User {
	return &User{
		Session: sess,
		Conn:    conn,
		outtime: time.Second * 10,
	}
}

type Agent struct {
	maps      map[uint]*User
	msg       func(*User, []byte) error
	leavefunc func(*User)
}

func NewAgent(max int, msg func(*User, []byte) error, le func(*User)) *Agent {
	return &Agent{
		msg:       msg,
		leavefunc: le,
		maps:      make(map[uint]*User, max),
	}
}

func (ag *Agent) Join(conn Conn) {
	obj := NewSession()
	user := NewUser(obj, conn)
	ag.loops(user)
}

//func (ag *Agent) JoinSync(conn Conn) *User {
//	obj := NewSession()
//	user := NewUser(obj, conn)
//	go ag.loops(user)
//	return user
//}

func (ag *Agent) loops(user *User) {
	uniq := user.Session.ToUint()
	ag.maps[uniq] = user
	base.NOTICE("Join ", user.RemoteAddr())
	ag.loop(user)
	if ag.leavefunc != nil {
		ag.leavefunc(user)
	}
	ag.leave(uniq)
	base.NOTICE("Leave ", user.RemoteAddr())
	user.Close()
}

func (ag *Agent) leave(uniq uint) {
	delete(ag.maps, uniq)
}

func (ag *Agent) Get(uniq uint) *User {
	return ag.maps[uniq]
}

func (ag *Agent) loop(user *User) {
	for {
		if b, err := user.Conn.ReadMsg(); err != nil {
			return
		} else {
			user.Lock()
			user.Session.refresh()
			err = ag.msg(user, b)
			user.Unlock()
			//			if err != nil {
			//				return
			//			}
		}
	}
}
