package main

import (
	"github.com/wzshiming/server"
	"github.com/wzshiming/server/agent"
	"github.com/wzshiming/server/agent/defaul"
	"github.com/wzshiming/server/cfg"
)

var ser *server.Server

func init() {
	cfg.TakeConf()
	ser = cfg.Self.Server()
}

func start() {

	ag := defaul.DefaulAgent()
	ser.Register(agent.NewConnect(ag))
	go agent.NewListener(cfg.Self.ClientPort, ag.Join).Start()
	ser.Start()
}

func main() {
	start()
}
