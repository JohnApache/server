package agent

import "github.com/wzshiming/base"

//"net/http"

type Request struct {
	Request *base.EncodeBytes
	Session *Session
	Head    []byte
}

type Response struct {
	Error    string
	Coverage *base.EncodeBytes
	Head     []byte
	Response *base.EncodeBytes
	//Data     *base.EncodeBytes

}

func (re Response) Hand(user *User, head []byte) error {
	var ret []byte
	if re.Coverage != nil {
		//base.INFO(string(re.Coverage.Bytes()))
		user.Data = re.Coverage
	}
	if re.Error != "" {
		ret = []byte(`{"error":"` + re.Error + `"}`)
	} else if re.Response != nil {
		ret = re.Response.Bytes()
	} else {
		return nil
		//ret = []byte(`{"error":""}`)
	}
	if re.Head != nil && len(re.Head) != 0 {
		ret = append(re.Head, ret...)
	} else {
		ret = append(head, ret...)
	}
	return user.WriteMsg(ret)
}
