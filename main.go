package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/VarthanV/kv-store/aof"
	"github.com/VarthanV/kv-store/command"
	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/resp"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
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
	defer conn.Close()

	commandClient := command.New()

	// Since Redis operates in a single threaded model we try to learn
	// and acheive the same here. Redis uses epoll and select to listen from multiple
	// fds simultaneously. This is just to make things simple
	for {
		request := resp.NewResp(conn)
		value, err := request.Read()
		if err != nil {
			logrus.Error(err)
			return
		}

		if value.Typ != objects.ARRAY {
			logrus.Error("invalid input expected input type ", objects.ARRAY)
			continue
		}

		if len(value.Arr) == 0 {
			logrus.Error("invalid length, expected args > 0")
			continue
		}

		writer := resp.NewWriter(conn)
		f, err := aof.New("aof.log")
		if err != nil {
			logrus.Error("persistence not enabled")
		}

		defer f.Close()

		cmd := strings.ToUpper(value.Arr[0].Bulk)
		args := value.Arr[1:]
		if commandClient.DoPersist(cmd) {
			f.Write(value)
		}

		result := commandClient.Handle(cmd, args)

		writer.Write(&result)

	}
}
