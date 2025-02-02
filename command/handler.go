package command

import (
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/pkg/utils"
	"github.com/VarthanV/kv-store/resp"
	"github.com/sirupsen/logrus"
)

type command string

const (
	ping command = "ping"
	get  command = "get"
	set  command = "set"
)

// HandlerFunc: Signature for the handler func to be implemented
type HandlerFunc func(args []resp.Value) resp.Value

// HandlerFuncMap: Map between command and handler func
type HandlerFuncMap map[string]HandlerFunc

type Handler struct {
	mu      sync.Mutex
	intsets map[string]*atomic.Int64
	sets    map[string]string
	hsets   map[string]map[string]string
}

func New() *Handler {
	return &Handler{
		mu:      sync.Mutex{},
		intsets: make(map[string]*atomic.Int64),
		sets:    make(map[string]string),
		hsets:   make(map[string]map[string]string),
	}
}

func (h *Handler) Handle(cmd string, args []resp.Value) resp.Value {
	switch command(cmd) {
	case ping:
		return h.ping(args)
	case set:
		return h.set(args)
	default:
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "Unknown command"}
	}
}

func (h *Handler) ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "PONG"}
	}
	return resp.Value{Typ: objects.SIMPLE_STRING, Str: args[0].Bulk}
}

func (h *Handler) set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: objects.ERROR_MESSAGE, Str: "Invalid number of arguments for SET"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	// If it is of type integer we can store in the
	// intset which is helpful in doing atomic INCR and DECR
	// operations
	if utils.IsInteger(value) {
		_val, err := strconv.Atoi(value)
		if err != nil {
			logrus.Error("error in converting to integer ", err)
			// Fallback to normal set
			h.mu.Lock()
			h.sets[key] = value
			h.mu.Unlock()
			logrus.Debugf("set key %s to value %s in normal set", key, value)
		} else {
			h.mu.Lock()
			var atomicVal atomic.Int64
			atomicVal.Store(int64(_val))
			h.intsets[key] = &atomicVal
			logrus.Debugf("set key %s to value %s in integer set", key, value)
		}
	}
	return resp.Value{Typ: objects.SIMPLE_STRING, Str: "OK"}
}
