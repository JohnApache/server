package agent

import "github.com/wzshiming/base"

//"net/http"

type Request struct {
	Head    []byte
	Request *base.EncodeBytes
	Session *Session
}

func (re Request) Mutex(reply *Response, f func()) {
	re.Session = re.Session.Sync()
	re.Session.Already(re, reply, f)
}

type Response struct {
	Head     []byte
	Response *base.EncodeBytes
	Session  *Session
}

func (re Response) Hand(user *User, head []byte) error {
	ret := []byte{}
	if re.Session != nil {
		user.Session = re.Session
	}

	if re.Response == nil {
		return nil
	}
	ret = re.Response.Bytes()

	if re.Head != nil && len(re.Head) != 0 {
		ret = append(re.Head, ret...)
	} else {
		ret = append(head, ret...)
	}
	return user.WriteMsg(ret)
}

func (re *Response) ReplyEncode(b *base.EncodeBytes) {
	re.Response = b
}

func (re *Response) Reply(s interface{}) {
	re.Response = base.EnJson(s)
}

func (re *Response) ReplyError(s interface{}) {
	re.Reply(map[string]interface{}{"error": s})
}
