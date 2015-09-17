package agent

import "github.com/wzshiming/base"

type rooms map[string]*Room

var ros = rooms{}

func SessionLeave(s *Session) {
	s.Mutex(func() {
		for _, v := range ros {
			v.Leave(s)
		}
		sesss.Leave(s)
	})
}

type Room struct {
	name string
	list map[uint]*Session
}

type data struct {
	Id   uint `json:",string"`
	Head []byte
}

type datafmt struct {
	Rooms map[string]data `json:"__Rooms__"`
}

func GetFromRooms(sess *Session) (r datafmt) {
	sess.Data.DeJson(&r)
	return
}

func GetFromRoom(sess *Session, name string) uint {
	return GetFromRooms(sess).Rooms[name].Id
}

func GetFromHead(sess *Session, name string) []byte {
	return GetFromRooms(sess).Rooms[name].Head
}

func SetFromHead(sess *Session, name string, head []byte) {
	var r datafmt
	sess.Data.DeJson(&r)
	r.Rooms[name] = data{
		Id:   r.Rooms[name].Id,
		Head: head,
	}
	sess.Data = base.SumJson(sess.Data, base.EnJson(r))
}

func NewRoom(name string) *Room {
	if ros[name] != nil {
		return nil
	}
	r := &Room{
		name: name,
		list: make(map[uint]*Session),
	}
	ros[name] = r
	return r
}

func (ro *Room) Repeal() {
	ro.ForEach(ro.Leave)
	ros[ro.name] = nil
}

func (ro *Room) JoinFrom(uniq uint, sess *Session, head []byte) {
	sess.Mutex(func() {
		var r datafmt
		sess.Data.DeJson(&r)
		if r.Rooms == nil {
			r.Rooms = make(map[string]data)
		}

		r.Rooms[ro.name] = data{
			Id:   uniq,
			Head: head,
		}
		ro.list[uniq] = sess
		ret := base.EnJson(r)
		sess.Data = base.SumJson(sess.Data, ret)
	})
	return
}

func (ro *Room) Join(sess *Session, head []byte) {
	ro.JoinFrom(sess.ToUint(), sess, head)
}

func (ro *Room) Leave(sess *Session) {
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

func (ro *Room) Uniq(sess *Session) uint {
	return GetFromRoom(sess, ro.name)
}

func (ro *Room) Head(sess *Session) []byte {
	return GetFromHead(sess, ro.name)
}

func (ro *Room) SetHead(sess *Session, head []byte) {
	SetFromHead(sess, ro.name, head)
}

func (ro *Room) Sync(sess *Session) *Session {
	if se := ro.Uniq(sess); se != 0 {
		return ro.Get(se)
	}
	return nil
}

func (ro *Room) Get(uniq uint) *Session {
	return ro.list[uniq]
}

func (ro *Room) Len() int {
	return len(ro.list)
}

func (ro *Room) ForEach(fun func(*Session)) {
	for _, v := range ro.list {
		fun(v)
	}
}

func (ro *Room) Group(name string, sesss ...*Session) (r *Room) {
	r = NewRoom(name)
	for _, v := range sesss {
		if i := ro.Sync(v); i != nil {
			head := ro.Head(i)
			ro.Leave(i)
			r.Join(i, head)
		}
	}
	return
}

func (ro *Room) GroupFromSize(size int) (sesss []*Session) {
	for _, v := range ro.list {
		if size == len(sesss) {
			return
		}
		sesss = append(sesss, v)
	}
	return
}

func (ro *Room) Push(reply interface{}, sess *Session) (err error) {
	return ro.Send(&Response{
		Response: base.EnJson(reply),
	},
		sess,
	)
}

func (ro *Room) Send(reply *Response, sess *Session) (err error) {
	if reply.Head == nil {
		reply.Head = ro.Head(sess)
	}
	if err = sess.Send(reply); err != nil {
		base.ERR(err)
		ro.Leave(sess)
	}
	return
}

func (ro *Room) BroadcastPush(reply interface{}, fail func(*Session)) {
	ro.Broadcast(&Response{
		Response: base.EnJson(reply),
	},
		fail,
	)
}

func (ro *Room) Broadcast(reply *Response, fail func(*Session)) {
	ro.ForEach(func(sess *Session) {
		if err := ro.Send(reply, sess); err != nil {
			if fail != nil {
				fail(sess)
			}
		}
	})
}
