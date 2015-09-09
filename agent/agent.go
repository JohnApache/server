package agent

import (
	"sync"
	"time"

	"github.com/wzshiming/base"
)

type User struct {
	Session
	Conn
	sync.RWMutex
	outtime time.Duration
}

func NewUser(sess *Session, conn Conn) *User {
	return &User{
		Session: *sess,
		Conn:    conn,
		outtime: time.Second * 10,
	}
}

type Agent struct {
	maps map[uint]*User
	msg  func(*User, []byte) error
}

func NewAgent(max int, msg func(*User, []byte) error) *Agent {
	return &Agent{
		msg:  msg,
		maps: make(map[uint]*User, max),
	}
}

func (ag *Agent) Join(conn Conn) {
	obj := NewSession()
	user := NewUser(obj, conn)
	ag.Loop(user)
}

func (ag *Agent) JoinSync(conn Conn) *User {
	obj := NewSession()
	user := NewUser(obj, conn)
	go ag.Loop(user)
	return user
}

func (ag *Agent) Loop(user *User) {
	uniq := user.Session.ToUint()
	ag.maps[uniq] = user
	base.NOTICE("Join ", user.RemoteAddr())
	ag.loop(user)
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
			err = ag.msg(user, b)
			user.Unlock()
			if err != nil {
				return
			}
		}
	}
}
