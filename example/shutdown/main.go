package main

import (
	"flag"
	"os"

	"github.com/wzshiming/server/cfg"
)

func init() {
	server := flag.String("server", "server.json", "")
	flag.Parse()
	if *server == "" {
		os.Exit(0)
	}

	cfg.Whole = cfg.NewWholeConfig(cfg.DirConf + *server)
}

func main() {
	defer func() {
		recover()
	}()
	cfg.Whole.Shutdown()
	cfg.Master.Client().ShutdownNow()
}
