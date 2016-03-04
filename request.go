/*
* @Author: anchen
* @Date:   2016-01-21 16:53:34
* @Last Modified by:   anchen
* @Last Modified time: 2016-01-21 17:15:01
 */

package redisChan

import (
	"bufio"
	"fmt"
	"strings"
)

// 1.-------------request------------------------

type reqBulk struct {
	// 请求命令
	cmd CommandCode
	// key -value
	key string
	// value
	value string
}

// 接收客户端的请求
// text protocol analysis
func (req *reqBulk) Receive(r *bufio.Reader) error {
	line, _, err := r.ReadLine()
	if err != nil || len(line) == 0 {
		return err
	}

	fmt.Println(string(line))

	params := strings.Fields(string(line))
	command := CommandCode(params[0])

	fmt.Println(command)
	req.cmd = command
	req.key = ""
	req.value = ""

	// keep it simple
	switch command {
	case GET, DELETE:
		req.cmd = command
		req.key = params[1]
	case STATS:
		req.cmd = command
		req.key = ""
	case SET:
		req.cmd = command
		req.key = params[1]
		req.value = params[2]
	}

	return err
}

// redis protocol analysis
func (req *reqBulk) RedisReceive(r *bufio.Reader) error {
	resp, err := Parse(r)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if resp == nil {
		fmt.Println("unknown error")
		return nil
	}

	cmd, err := resp.Op()
	if err != nil {
		fmt.Println("can not get cmd")
		return err
	}
	req.cmd = CommandCode(cmd)

	key, err := resp.Key()
	if err != nil {
		req.key = ""
	} else {
		req.key = string(key)
	}
	value, err := resp.Value()
	if err != nil {
		req.value = ""
	} else {
		req.value = string(value)
	}

	fmt.Println(req)

	return nil

}

func (req *reqBulk) String() string {
	return fmt.Sprintf("{Request opcode=%s, key=%s, value='%s'}", req.cmd, req.key, req.value)
}
