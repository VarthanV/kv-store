package main

import (
	"fmt"
	"net"

	"github.com/VarthanV/kv-store/resp"
	"github.com/sirupsen/logrus"
)

func main() {
	const PORT = "6363"
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		logrus.Error("error in listening ", err)
		return

	}

	logrus.Info("Listening on port ", PORT)

	conn, err := l.Accept()
	if err != nil {
		logrus.Error("error in accepting connection ", err)
		return
	}

	// Since Redis operates in a single threaded model we try to learn
	// and acheive the same here. Redis uses epoll and select to listen from multiple
	// fds simultaneously. This is just to make things simple
	for {
		resp := resp.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", value)

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
		// PONG
	}

	defer conn.Close()

}
