package cfg

import (
	"flag"
)

var (
	DirConf   = "../conf/"
	DirI18n   = "../i18n/"
	DirDbs    = "../dbs/"
	DirPublic = "../public/"
	DirView   = "../view/"
	DirScript = "../script/"
	DirTemp   = "../tmep/"
)

var (
	Master *ServerConfig
	Whole  *WholeConfig
	Self   *ServerConfig
	SelfId int
)

func TakeConf() {
	Master = NewServerConfig(DirConf + "master.json")
	takeWhole()
	takeSelf()
}

func takeSelf() {
	if Self != nil {
		return
	}
	flag.IntVar(&SelfId, "id", -1, "server unique id ")
	flag.Parse()
	Self = ConfForId[SelfId]
}

func takeWhole() {
	if Whole != nil {
		return
	}
	var w WholeConfig
	err := Master.Client().Take("Master.WholeConfig", &w)
	if err != nil {
		return
	}

	for k, v := range w.Apps {
		ConfForId[v.Id] = &w.Apps[k]
	}
	for k, v := range w.Agents {
		ConfForId[v.Id] = &w.Agents[k]
	}
	Whole = &w
}
