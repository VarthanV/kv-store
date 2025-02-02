package command

import (
	"sync"
	"sync/atomic"

	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/resp"
)

// HandlerFunc: Signature for the handler func to be implemented
type HandlerFunc func(args []resp.Value) resp.Value

// HandlerFuncMap: Map between command and handler func
type HandlerFuncMap map[string]HandlerFunc

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "PONG"}
	}
	return resp.Value{Typ: objects.SIMPLE_STRING, Str: args[0].Bulk}
}

var handlers = HandlerFuncMap{
	"ping": ping,
}

type Handler struct {
	mu      sync.Mutex
	Map     HandlerFuncMap
	intsets map[string]atomic.Int64
	sets    map[string]string
	hsets   map[string]map[string]string
}

func New() *Handler {
	return &Handler{
		Map:     handlers,
		mu:      sync.Mutex{},
		intsets: make(map[string]atomic.Int64),
		sets:    make(map[string]string),
		hsets:   make(map[string]map[string]string),
	}
}

func (h *Handler) GetHandler(command string) (HandlerFunc, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	handler, ok := h.Map[command]
	return handler, ok
}
