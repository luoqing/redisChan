package redisChan

import (
	"fmt"
	"io"
)

//status result map
var StatusRes map[Status]string

//初始化response状态的返回结果
func init() {
	StatusRes = make(map[Status]string)
	StatusRes[ERROR] = "ERROR"
	StatusRes[STORED] = "STORED"
	StatusRes[NOT_STORED] = "NOT_STORED"
	StatusRes[END] = "END"
	StatusRes[DELETED] = "DELETED"
	StatusRes[NOT_FOUND] = "NOT_FOUND"
}

//status to string
func (s *Status) ToString() string {
	rv := StatusRes[*s]
	if rv == "" {
		rv = fmt.Sprintf("%s", StatusRes[NOT_FOUND])
	} else {
		rv = fmt.Sprintf("%s", rv)
	}
	return rv
}

type responseBulk struct {
	cmd    CommandCode
	key    string
	value  string
	status Status
}

//解析response 并把返回结果写入socket链接
func (res *responseBulk) Transmit(w io.Writer) (err error) {
	var response string
	switch res.cmd {
	case GET:
		if res.status == SUCCESS {
			response = res.value
		} else {
			response = res.status.ToString()
		}
	case DELETE:
		response = "DELETED"
	default:
		response = res.status.ToString()
	}

	response = fmt.Sprintf("%s\r\n", response)

	_, err = w.Write([]byte(response))

	return err
}

func (res *responseBulk) RedisTransmit(w io.Writer) (err error) {
	var response string
	switch res.cmd {
	case GET:
		if res.status == SUCCESS {
			response = res.value
		} else {
			response = res.status.ToString()
		}
	case DELETE:
		response = "DELETED"
	default:
		response = res.status.ToString()
	}
	resp := &RESP{
		t:     SimpleString,
		b:     []byte(response),
		array: nil,
	}

	b, err := resp.Bytes()
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
