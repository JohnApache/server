package route

import (
	"github.com/wzshiming/base"
	"github.com/wzshiming/server"
	"github.com/wzshiming/server/cfg"
)

type Route struct {
	maps     *CodeMaps
	connmaps map[string][]*server.Client
	conns    []*server.Client
}

func NewRoute(conf []cfg.ServerConfig) *Route {
	ro := &Route{
		connmaps: make(map[string][]*server.Client),
		maps:     NewCodeMaps(),
	}
	for _, v1 := range conf {
		ro.Register(v1)
	}
	return ro
}

func (ro *Route) Code() *CodeMaps {
	return ro.maps
}

func (ro *Route) Register(c cfg.ServerConfig) {
	conn := c.Client()
	if conn == nil {
		return
	}

	ro.conns = append(ro.conns, conn)
	if len(ro.connmaps[c.Type]) == 0 && c.Type != "" {
		classs, err := conn.TakeClasss()
		if err != nil {
			return
		}
		ro.maps.Append(c.Type, classs)
	}
	ro.connmaps[c.Type] = append(ro.connmaps[c.Type], conn)
}

func (ro *Route) CallCode(c1, c2, c3 byte, args, reply interface{}) error {
	m1, m2, m3, err := ro.maps.Map(c1, c2, c3)
	if err != nil {
		return err
	}
	return ro.Call(m1, m2, m3, args, reply)
}

func (ro *Route) Call(m1, m2, m3 string, args, reply interface{}) (err error) {
	defer base.PanicErr(&err, "Route.Call: index out of range")

	return ro.connmaps[m1][0].Call(m2+"."+m3, args, reply)
}
