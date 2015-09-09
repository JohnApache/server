package defaul

import (
	"errors"
	"fmt"
	"time"

	rego "github.com/wzshiming/base"
	"github.com/wzshiming/server/agent"
	"github.com/wzshiming/server/cfg"
	"github.com/wzshiming/server/route"
)

var MapFile = "map.json"

func DefaulAgent() *agent.Agent {
	ro := route.NewRoute(cfg.Whole.Apps)
	ro.Code().WriteFile(MapFile)
	ag := agent.NewAgent(1024, func(user *agent.User, msg []byte) (err error) {
		defer rego.PanicErr(&err)
		var reply agent.Response
		//rego.INFO(string(msg))
		user.SetDeadline(time.Now().Add(time.Second * 60 * 60))
		user.Refresh()
		err = ro.CallCode(msg[1], msg[2], msg[3], agent.Request{
			Session: &user.Session,
			Request: rego.NewEncodeBytes(msg[4:]),
			Head:    msg[:4],
		}, &reply)

		if err != nil {
			ret := []byte(`{"error":"` + err.Error() + `"}`)
			return user.WriteMsg(append(msg[:4], ret...))
		}
		return reply.Hand(user, msg[:4])

	})
	return ag
}

func DefaulConn() agent.Conn {
	ca := cfg.Whole.Agents[0].ClientAddr()
	return agent.NewConn(ca)
}

func DefaulClient(hand func(byte, byte, byte, []byte) error) func(byte, byte, byte, []byte) error {
	conn := DefaulConn()
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

func DefaultClientCode(hand func(code string, v interface{}) error) func(code string, v interface{}) error {
	cod := route.NewCodeMaps()
	cod.ReadFile(MapFile)
	recod := cod.MakeReCodeMap()
	rego.INFO(recod)
	c := DefaulClient(func(m1 byte, m2 byte, m3 byte, b []byte) error {
		c1, c2, c3, err := cod.Map(m1, m2, m3)
		var code string
		if err != nil {
			code = fmt.Sprintf("None.%d.%d.%d", m1, m2, m3)
		} else {
			code = c1 + "." + c2 + "." + c3
		}
		es := rego.NewEncodeBytes(b)
		var r interface{}
		es.DeJson(&r)
		return hand(code, r)
	})
	return func(code string, v interface{}) error {
		m1, m2, m3, err := recod.Map(code)
		if err != nil {
			rego.ERR(err)
			return err
		}
		es := rego.EnJson(v)
		return c(m1, m2, m3, es.Bytes())
	}
}
