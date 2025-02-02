package command

import (
	"sync"
	"sync/atomic"

	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/resp"
)

type command string

const (
	ping command = "ping"
)

// HandlerFunc: Signature for the handler func to be implemented
type HandlerFunc func(args []resp.Value) resp.Value

// HandlerFuncMap: Map between command and handler func
type HandlerFuncMap map[string]HandlerFunc

type Handler struct {
	mu      sync.Mutex
	intsets map[string]atomic.Int64
	sets    map[string]string
	hsets   map[string]map[string]string
}

func New() *Handler {
	return &Handler{
		mu:      sync.Mutex{},
		intsets: make(map[string]atomic.Int64),
		sets:    make(map[string]string),
		hsets:   make(map[string]map[string]string),
	}
}

func (h *Handler) Handle(cmd string, args []resp.Value) resp.Value {
	switch command(cmd) {
	case ping:
		return h.ping(args)
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
