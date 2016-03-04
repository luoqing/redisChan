package redisChan

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// 0.-------------listen------------------------
//一个 请求产生 一个 chanReq
type chanReq struct {
	request  *reqBulk
	response chan *responseBulk
}

//请求句柄
type reqHandler struct {
	ch chan chanReq
}

//
func handleConn(c net.Conn, handler *reqHandler) error {
	defer c.Close()
	for {
		r := bufio.NewReader(c)
		req := &reqBulk{}
		// err := req.Receive(r) // read
		err := req.RedisReceive(r) // read

		if err != nil {
			fmt.Println("not receive")
			return err
		}

		cr := chanReq{
			req,
			make(chan *responseBulk),
		}

		handler.ch <- cr     // 有请求写过来
		res := <-cr.response // 将处理结果返回过来

		// res.Transmit(c) // write
		res.RedisTransmit(c) // write
		if err != nil {
			return err
		}
	}
	return nil
}

func runServer() {
	// listen
	port := 12354
	listener, error := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if error != nil {
		log.Fatalf("Got an ERROR: %s", error)
	}
	fmt.Println("start listening port:", port)
	reqChannel := make(chan chanReq)
	go waitDispatch(reqChannel) // dispatch to handler
	handler := &reqHandler{reqChannel}
	// accept
	for {
		conn, error := listener.Accept()
		if error != nil {
			log.Printf("Got an ERROR:%s", error)
			continue
		} else {
			fmt.Println("get a connection")
			go handleConn(conn, handler)
		}
	}
}
