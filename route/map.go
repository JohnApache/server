package route

import (
	"errors"
	"io/ioutil"

	"github.com/wzshiming/base"
	"github.com/wzshiming/server"
)

type CodeMap struct {
	Name   string
	Classs server.Classs
}

type CodeMaps []CodeMap

func NewCodeMaps() *CodeMaps {
	r := CodeMaps{}
	return &r
}

func (co *CodeMaps) Map(c1, c2, c3 byte) (m1, m2, m3 string, err error) {

	defer base.PanicErr(&err, "CodeMaps.Map: index out of range")
	i1 := (*co)[c1]
	i2 := i1.Classs[c2]
	m3 = i2.Methods[c3]
	m2 = i2.Name
	m1 = i1.Name
	return
}

func (co *CodeMaps) Append(name string, classs server.Classs) {
	*co = append(*co, CodeMap{
		Name:   name,
		Classs: classs,
	})
}

func (co *CodeMaps) MakeReCodeMap() *ReCodeMaps {
	re := make(ReCodeMaps)
	for k1, v1 := range *co {
		for k2, v2 := range v1.Classs {
			for k3, v3 := range v2.Methods {
				re[v1.Name+"."+v2.Name+"."+v3] = K2i(k1, k2, k3)
			}
		}
	}
	return &re
}

func (co *CodeMaps) WriteFile(name string) {
	ioutil.WriteFile(name, base.EnJson(co).Bytes(), 0666)
}

func (co *CodeMaps) ReadFile(name string) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		base.ERR(err)
		return
	}
	base.NewEncodeBytes(b).DeJson(co)
}

type ReCodeMaps map[string]uint32

func (co *ReCodeMaps) Register(code string, k1, k2, k3 int) {
	(*co)[code] = K2i(k1, k2, k3)
}

func (co *ReCodeMaps) Map(code string) (byte, byte, byte, error) {
	return I2k((*co)[code])
}

func K2i(k1, k2, k3 int) uint32 {
	return uint32(k1 | (k2 << 8) | (k3 << 16) | (1 << 24))
}

func I2k(i uint32) (byte, byte, byte, error) {
	if (uint32(i) >> 24) != 1 {
		return 255, 255, 255, errors.New("I2k: inexistence")
	}
	return byte(i), byte(i >> 8), byte(i >> 16), nil
}
