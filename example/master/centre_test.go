package main

import (
	"testing"

	"github.com/wzshiming/server/cfg"
)

func Test_get(t *testing.T) {
	go start()
	cfg.TakeConf()
}
