package aof

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/VarthanV/kv-store/resp"
	"github.com/sirupsen/logrus"
)

type Aof struct {
	f  *os.File
	rd *bufio.Reader
	mu sync.Mutex
}

func New(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		f:  f,
		rd: bufio.NewReader(f),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 30 * 1) // every 2 mins
		for range ticker.C {
			logrus.Info("Flush routine started")
			aof.mu.Lock()
			err := aof.f.Sync()
			if err != nil {
				logrus.Error("error in flushing file ", err)
			}
			aof.mu.Unlock()
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.f.Close()
}

func (aof *Aof) Write(value *resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	logrus.Info("writing to aof")
	_, err := aof.f.Write(value.Marshal())
	if err != nil {
		return err
	}
	aof.f.Sync()
	return nil
}

func (a *Aof) Read(callback func(val *resp.Value)) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	resp := resp.NewResp(a.f)

	for {
		value, err := resp.Read()
		if err == nil {
			callback(value)
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
