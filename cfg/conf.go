package cfg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server"
)

var index = 0
var ConfForId = make(map[int]*ServerConfig)

type ServerConfig struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	ClientPort int    `json:"clientport"`
	Bin        string `json:"bin"`
	Id         int    `json:"id"`
}

func (se *ServerConfig) makeId() {
	for ConfForId[index] != nil {
		index++
	}
	se.Id = index
	ConfForId[index] = se

	if se.Name == "" {
		se.Name = fmt.Sprintf("%s_%s_%d", se.Type, se.Bin, se.Id)
	}
}

func (se *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", se.Host, se.Port)
}

func (se *ServerConfig) ClientAddr() string {
	return fmt.Sprintf("%s:%d", se.Host, se.ClientPort)
}

func (se *ServerConfig) ExposedAddr() string {
	if se.ClientPort != 0 {
		return fmt.Sprintf("%s:%d", base.LocalIP, se.Port)
	}
	return ""
}

func (se *ServerConfig) Client() *server.Client {
	return server.NewClient(se.Addr())
}

func (se *ServerConfig) Server() *server.Server {
	return server.NewServer(se.Port)
}

func (se *ServerConfig) Start() {
	b := fmt.Sprintf("./%s", strings.ToLower(se.Type))
	args := []string{}
	if se.Bin != "" {
		b = fmt.Sprintf("./%s", se.Bin)
	}
	args = append(args, "-id", fmt.Sprintf("%d", se.Id))
	cmd := exec.Command(b, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	base.NOTICE("Start", se.Id, se.Type, b)
	err := cmd.Start()
	if err != nil {
		base.ERR(err)
	}
}
