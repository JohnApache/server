package server

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

var (
	request  = "Request"
	response = "Response"
)

type Methods struct {
	Name    string
	Methods []string
}

func NewMethods(rcvr interface{}) *Methods {
	an := Methods{}
	typ := reflect.TypeOf(rcvr)
	an.Name = reflect.Indirect(reflect.ValueOf(rcvr)).Type().Name()

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name

		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			continue
		}
		argType := mtype.In(1)
		if argType.Name() != request {
			continue
		}
		// First arg need not be a pointer.
		if !isExportedOrBuiltinType(argType) {
			continue
		}

		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			continue
		}

		if replyType.Elem().Name() != response {
			continue
		}

		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			continue
		}

		// Method needs one out.
		if mtype.NumOut() != 1 {
			continue
		}

		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}
		an.Methods = append(an.Methods, mname)
	}
	return &an
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

type Classs []Methods

func (cl *Classs) Take(args int, reply *Classs) error {
	*reply = *cl
	return nil
}
