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

	cfg.Master = cfg.NewServerConfig(cfg.DirConf + "master.json")
	cfg.Whole = cfg.NewWholeConfig(cfg.DirConf + *server)
	cfg.Self = cfg.Master
}

type Master struct {
}

func NewMaster() *Master {
	return &Master{}
}

func (ma *Master) WholeConfig(args int, reply *cfg.WholeConfig) error {
	*reply = *cfg.Whole
	return nil
}

func (ma *Master) Shutdown(int, *int) error {
	cfg.Whole.Shutdown()
	return nil
}

func start() {
	ser := cfg.Self.Server()
	ser.Register(NewMaster())

	go cfg.Whole.Start()
	ser.Start()
}

func main() {
	start()
}
