package main

import (
	"fmt"
	"io"
	"net"
	"os"

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
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error("error reading from client: ", err)
			os.Exit(1)
		}

		// PONG
		conn.Write([]byte("+OK\r\n"))
	}

	defer conn.Close()

}
