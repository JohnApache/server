package main

import (
	"github.com/wzshiming/server/cfg"
)

func init() {
	cfg.Whole = cfg.NewWholeConfig(cfg.DirConf + "server.json")
}

func main() {
	defer func() {
		recover()
	}()
	cfg.Whole.Shutdown()
	cfg.Master.Client().ShutdownNow()
}
