package rooms

import (
	"errors"
	//"fmt"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server/agent"
)

type Rooms struct {
	name string
	list map[uint]*agent.Session
}

type data struct {
	ID   uint `json:",string"`
	Head []byte
}

type datafmt struct {
	Rooms map[string]data `json:"__Rooms__"`
}

func GetFromRoom(sess *agent.Session, name string) uint {
	var r datafmt
	sess.Data.DeJson(&r)
	return r.Rooms[name].ID
}

func GetFromHead(sess *agent.Session, name string) []byte {
	var r datafmt
	sess.Data.DeJson(&r)
	return r.Rooms[name].Head
}

func SetFromHead(sess *agent.Session, name string, head []byte) {
	var r datafmt
	sess.Data.DeJson(&r)
	r.Rooms[name] = data{
		ID:   r.Rooms[name].ID,
		Head: head,
	}
	sess.Data = base.SumJson(sess.Data, base.EnJson(r))
}

func NewRooms(name string) *Rooms {
	return &Rooms{
		name: name,
		list: make(map[uint]*agent.Session),
	}
}

func (ro *Rooms) JoinFrom(uniq uint, sess *agent.Session, head []byte) (err error) {
	sess.Mutex(func() {
		var r datafmt
		sess.Data.DeJson(&r)
		if r.Rooms == nil {
			r.Rooms = make(map[string]data)
		}
		if se := r.Rooms[ro.name].ID; se != 0 {
			err = errors.New("Rooms.Join: already join")
			return
		}
		if un := ro.list[uniq]; un != nil {
			err = errors.New("Rooms.Join: repetition join")
			return
		}
		r.Rooms[ro.name] = data{
			ID:   uniq,
			Head: head,
		}
		ro.list[uniq] = sess
		ret := base.EnJson(r)
		sess.Data = base.SumJson(sess.Data, ret)
	})
	return
}

func (ro *Rooms) Join(sess *agent.Session, head []byte) error {
	return ro.JoinFrom(sess.ToUint(), sess, head)
}

func (ro *Rooms) Leave(sess *agent.Session) {
	sess.Mutex(func() {
		uniq := ro.Uniq(sess)
		var r datafmt
		sess.Data.DeJson(&r)
		delete(r.Rooms, ro.name)
		delete(ro.list, uniq)
		re := base.EnJson(r)
		sess.Data = base.SumJson(sess.Data, re)
	})
	return
}

func (ro *Rooms) Uniq(sess *agent.Session) uint {
	return GetFromRoom(sess, ro.name)
}

func (ro *Rooms) Head(sess *agent.Session) []byte {
	return GetFromHead(sess, ro.name)
}

func (ro *Rooms) SetHead(sess *agent.Session, head []byte) {
	SetFromHead(sess, ro.name, head)
}

func (ro *Rooms) Sync(sess *agent.Session) *agent.Session {
	if se := ro.Uniq(sess); se != 0 {
		return ro.Get(se)
	}
	return nil
}

func (ro *Rooms) Get(uniq uint) *agent.Session {
	return ro.list[uniq]
}

func (ro *Rooms) Len() int {
	return len(ro.list)
}

func (ro *Rooms) ForEach(fun func(*agent.Session)) {
	for _, v := range ro.list {
		fun(v)
	}
}

func (ro *Rooms) Group(name string, sesss ...*agent.Session) (r *Rooms) {
	r = NewRooms(name)
	for _, v := range sesss {
		if i := ro.Sync(v); i != nil {
			head := ro.Head(i)
			ro.Leave(i)
			r.Join(i, head)
		}
	}
	return
}

func (ro *Rooms) GroupFromSize(size int) (sesss []*agent.Session) {
	for _, v := range ro.list {
		if size == len(sesss) {
			return
		}
		sesss = append(sesss, v)
	}
	return
}

func (ro *Rooms) Push(reply interface{}, sess *agent.Session) (err error) {
	return ro.Send(&agent.Response{
		Response: base.EnJson(reply),
	},
		sess,
	)
}

func (ro *Rooms) Send(reply *agent.Response, sess *agent.Session) (err error) {
	if reply.Head == nil {
		reply.Head = ro.Head(sess)
	}
	if err = sess.Send(reply); err != nil {
		base.ERR(err)
		ro.Leave(sess)
	}
	return
}

func (ro *Rooms) BroadcastPush(reply interface{}, fail func(*agent.Session)) {
	ro.Broadcast(&agent.Response{
		Response: base.EnJson(reply),
	},
		fail,
	)
}

func (ro *Rooms) Broadcast(reply *agent.Response, fail func(*agent.Session)) {
	ro.ForEach(func(sess *agent.Session) {
		if err := ro.Send(reply, sess); err != nil {
			if fail != nil {
				fail(sess)
			}
		}
	})
}
