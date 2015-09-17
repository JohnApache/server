package defaul

import (
	"errors"
	"fmt"
	"time"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server/agent"
	"github.com/wzshiming/server/cfg"
	"github.com/wzshiming/server/route"
)

//var MapFile = "map.json"

type DefCfg struct {
	CodeMaps *route.CodeMaps
	Agents   []cfg.ServerConfig
}

func DefaulAgent() *agent.Agent {
	ro := route.NewRoute(cfg.Whole.Apps)
	code := ro.Code()
	recode := code.MakeReCodeMap()
	dj := base.EnJson(DefCfg{
		CodeMaps: code,
		Agents:   cfg.Whole.Agents,
	}).Bytes()
	base.INFO(recode)
	ag := agent.NewAgent(1024, func(user *agent.User, msg []byte) (err error) {

		defer base.PanicErr(&err)
		var reply agent.Response
		user.SetDeadline(time.Now().Add(time.Second * 60 * 60))
		user.Refresh()
		if msg[0] == 0 && msg[1] == 0 && msg[2] == 0 && msg[3] == 0 {
			return user.WriteMsg(append(msg[:4], dj...))
		}
		err = ro.CallCode(msg[1], msg[2], msg[3], agent.Request{
			Session: &user.Session,
			Request: base.NewEncodeBytes(msg[4:]),
			Head:    msg[:4],
		}, &reply)
		//		if err != nil {
		//			return nil
		//			ret := []byte(`{"error":"` + err.Error() + `"}`)
		//			return user.WriteMsg(append(msg[:4], ret...))
		//		}

		return reply.Hand(user, msg[:4])
	}, func(user *agent.User) {
		for _, v := range *code {
			sess := &user.Session
			ro.Call(v.Name, "Leave", "User", sess, nil)
		}
		//		sess := &user.Session
		//		rooms := agent.GetFromRooms(sess)
		//		for k, v := range rooms {

		//		}
		//recode
	})
	return ag
}

func DefaulConn() agent.Conn {
	ca := cfg.Whole.Agents[0].ClientAddr()
	return agent.NewConn(ca)
}

func DefaulClient(addr string, hand func(byte, byte, byte, []byte) error) func(byte, byte, byte, []byte) error {
	conn := agent.NewConn(addr)
	size := 0
	isEnd := false
	errmsg := errors.New("use of closed network connection")
	go func() {
		for {
			b, err := conn.ReadMsg()
			if err != nil {
				break
			}
			if len(b) == 0 {
				continue
			}
			err = hand(b[1], b[2], b[3], b[4:])
			if err != nil {
				break
			}
		}
		isEnd = true
		conn.Close()
	}()
	return func(m1, m2, m3 byte, b []byte) error {
		if isEnd {
			return errmsg
		}
		err := conn.WriteMsg(append([]byte{byte(size), m1, m2, m3}, b...))
		if err != nil {
			isEnd = true
			return errmsg
		}
		return nil
	}
}

func DefaultClientCode(addr string, hand func(code string, v interface{}) error) func(code string, v interface{}) error {

	ro := route.NewRoute(cfg.Whole.Apps)
	cod := ro.Code()
	recod := cod.MakeReCodeMap()
	base.INFO(recod)
	c := DefaulClient(addr, func(m1 byte, m2 byte, m3 byte, b []byte) error {
		c1, c2, c3, err := cod.Map(m1, m2, m3)
		var code string
		if err != nil {
			code = fmt.Sprintf("None.%d.%d.%d", m1, m2, m3)
		} else {
			code = c1 + "." + c2 + "." + c3
		}
		es := base.NewEncodeBytes(b)
		var r interface{}
		es.DeJson(&r)
		return hand(code, r)
	})
	return func(code string, v interface{}) error {
		m1, m2, m3, err := recod.Map(code)
		if err != nil {
			base.ERR(err)
			return err
		}
		es := base.EnJson(v)
		return c(m1, m2, m3, es.Bytes())
	}
}
