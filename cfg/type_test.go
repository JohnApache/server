package cfg

import (
	"testing"

	"github.com/wzshiming/base"
)

func Test_read(t *testing.T) {
	b := NewWholeConfig("./conf/server.json")
	if b == nil {
		t.Fail()
	}
	c := NewServerConfig("./conf/master.json")
	if c == nil {
		t.Fail()
	}
	base.INFO(b)
	base.INFO(c)
}
