package route

import (
	"testing"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server"
)

func Test_code(t *testing.T) {
	cm := NewCodeMaps()
	cm.Append("he", server.Classs{
		server.Methods{
			Name: "ll",
			Methods: []string{
				"oo",
				"o",
			},
		},
	})
	base.INFO(cm)
	base.INFO(cm.MakeReCodeMap().Map("he.ll.o"))
}
