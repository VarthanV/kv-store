package command

import (
	"sync"

	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/resp"
	"github.com/sirupsen/logrus"
)

type command string

const (
	ping command = "PING"
	get  command = "GET"
	set  command = "SET"
)

// HandlerFunc: Signature for the handler func to be implemented
type HandlerFunc func(args []resp.Value) resp.Value

// HandlerFuncMap: Map between command and handler func
type HandlerFuncMap map[string]HandlerFunc

type Handler struct {
	mu          sync.Mutex
	sets        map[string]string
	hsets       map[string]map[string]string
	keyMetaData map[string]string
}

func New() *Handler {
	return &Handler{
		mu:          sync.Mutex{},
		sets:        make(map[string]string),
		hsets:       make(map[string]map[string]string),
		keyMetaData: map[string]string{},
	}
}

func (h *Handler) Handle(cmd string, args []resp.Value) resp.Value {
	switch command(cmd) {
	case ping:
		return h.ping(args)
	case set:
		return h.set(args)
	case get:
		return h.get(args)
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
		logrus.Error("invalid number of arguments")
		return resp.Value{Typ: objects.ERROR_MESSAGE, Str: "Invalid number of arguments for SET"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	h.mu.Lock()
	h.sets[key] = value
	h.mu.Unlock()
	logrus.Debugf("set key %s to value %s in set", key, value)

	return resp.Value{Typ: objects.SIMPLE_STRING, Str: "OK"}
}

func (h *Handler) get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		logrus.Error("invalid number of arguments")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "Invalid number of arguments for GET"}
	}

	key := args[0]

	h.mu.Lock()
	defer h.mu.Unlock()

	val, ok := h.sets[key.Bulk]
	if ok {
		h.mu.Unlock()
		return resp.Value{Typ: objects.BULK_STRING, Bulk: val}
	}

	return resp.Value{Typ: objects.NULL}
}
