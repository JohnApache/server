package agent

type rooms map[string]*Room

var ros = rooms{}

func sessionLeave(s *Session) {
	s = s.Sync()
	s.NonSync(func() {
		for k, _ := range s.Rooms.Data() {
			if ros[k] != nil {
				ros[k].Leave(s)
			}
			//			else {
			//				base.ERR(ros)
			//			}
		}

		sesss.Leave(s)
	})
}

type Room struct {
	name       string
	user       map[uint]*Session
	parent     *Room
	child      map[string]*Room
	joinEvent  func(sess *Session)
	leaveEvent func(sess *Session)
}

type data struct {
	Id   uint `json:",string"`
	Head []byte
}

type roomsData map[string]data

func newRoomChild(name string, parent *Room) *Room {
	if parent != nil {
		name = parent.name + "." + name
	}
	if ros[name] != nil {
		return nil
	}
	r := &Room{
		parent: parent,
		name:   name,
		user:   map[uint]*Session{},
		child:  map[string]*Room{},
	}
	ros[name] = r
	return r
}

func NewRoom(name string) *Room {
	return newRoomChild(name, nil)
}

func (ro *Room) Name() string {
	return ro.name
}

func (ro *Room) Close() {
	ro.ForEach(func(sess *Session) {
		ro.Leave(sess)
	})
	ros[ro.name] = nil
}

func (ro *Room) JoinFrom(uniq uint, sess *Session, head []byte) {
	sess.Rooms.Set(ro.name, data{
		Id:   uniq,
		Head: head,
	})
	ro.user[uniq] = sess
	if ro.joinEvent != nil {
		ro.joinEvent(sess)
	}
	return
}

func (ro *Room) Join(sess *Session, head []byte) {
	ro.JoinFrom(sess.ToUint(), sess, head)
}

func (ro *Room) LeaveFrom(uniq uint) (d data) {
	sess := ro.user[uniq]
	if sess == nil {
		return
	}
	sess.Rooms.Get(ro.name, &d)
	sess.Rooms.Del(ro.name)
	delete(ro.user, uniq)
	if ro.leaveEvent != nil {
		ro.leaveEvent(sess)
	}
	return
}

func (ro *Room) Leave(sess *Session) data {
	return ro.LeaveFrom(ro.Uniq(sess))
}

func (ro *Room) GetChild(name string) (nr *Room) {
	nr = ro.child[name]
	if nr == nil {
		nr = newRoomChild(name, ro)
		ro.child[name] = nr
	}
	return
}

func (ro *Room) ToChild(sess *Session, name string) (nr *Room) {
	d := ro.Leave(sess)
	if d.Id != 0 {
		nr = ro.GetChild(name)
		nr.JoinFrom(d.Id, sess, d.Head)
	}
	return
}

func (ro *Room) GetParent() *Room {
	return ro.parent
}

func (ro *Room) ToParent(sess *Session) (nr *Room) {
	d := ro.Leave(sess)
	if nr = ro.parent; nr != nil && d.Id != 0 {
		nr.JoinFrom(d.Id, sess, d.Head)
	}
	return
}

func (ro *Room) Uniq(sess *Session) uint {
	return sess.RoomsUniq(ro.name)
}

func (ro *Room) Head(sess *Session) []byte {
	return sess.RoomsHead(ro.name)
}

func (ro *Room) SetHead(sess *Session, head []byte) {
	d := data{}
	sess.Rooms.Get(ro.name, &d)
	d.Head = head
	sess.Rooms.Set(ro.name, d)
}

func (ro *Room) IsExist(sess *Session) bool {
	if se := ro.Uniq(sess); se != 0 {
		return true
	}
	return false
}

func (ro *Room) Sync(sess *Session) *Session {
	if se := ro.Uniq(sess); se != 0 {
		return ro.Get(se)
	}
	return nil
}

func (ro *Room) Get(uniq uint) *Session {
	return ro.user[uniq]
}

func (ro *Room) Len() int {
	return len(ro.user)
}

func (ro *Room) ForEach(fun func(*Session)) {
	for _, v := range ro.user {
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
	for _, v := range ro.user {
		if size == len(sesss) {
			return
		}
		sesss = append(sesss, v)
	}
	return
}

func (ro *Room) Push(reply interface{}, sess *Session) (err error) {
	return sess.Push(reply, ro.Head(sess))
}

func (ro *Room) Broadcast(reply interface{}) {
	ro.ForEach(func(sess *Session) {
		ro.Push(reply, sess)
	})
}

func (ro *Room) JoinEvent(f func(sess *Session)) {
	ro.joinEvent = f
}

func (ro *Room) LeaveEvent(f func(sess *Session)) {
	ro.leaveEvent = f
}
