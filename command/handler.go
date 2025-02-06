package command

import (
	"strconv"
	"sync"

	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/pkg/utils"
	"github.com/VarthanV/kv-store/resp"
	"github.com/sirupsen/logrus"
)

type command string

const (
	ping    command = "PING"
	get     command = "GET"
	set     command = "SET"
	hset    command = "HSET"
	hget    command = "HGET"
	hgetall command = "HGETALL"
	del     command = "DEL"
	incr    command = "INCR"
	decr    command = "DECR"
	rappend command = "APPEND"
	lpush   command = "LPUSH"
)

// HandlerFunc: Signature for the handler func to be implemented
type HandlerFunc func(args []resp.Value) resp.Value

// HandlerFuncMap: Map between command and handler func
type HandlerFuncMap map[string]HandlerFunc

type Handler struct {
	mu          sync.Mutex
	sets        map[string]string
	hsets       map[string]map[string]string
	lists       map[string][]string
	keyMetaData map[string]string
}

func New() *Handler {
	return &Handler{
		mu:          sync.Mutex{},
		sets:        make(map[string]string),
		hsets:       make(map[string]map[string]string),
		keyMetaData: map[string]string{},
		lists:       make(map[string][]string),
	}
}

func okResponse() resp.Value {
	return resp.Value{Typ: objects.SIMPLE_STRING, Str: "OK"}
}

func (h *Handler) Handle(cmd string, args []resp.Value) resp.Value {
	switch command(cmd) {
	case ping:
		return h.ping(args)
	case set:
		return h.set(args)
	case get:
		return h.get(args)
	case hset:
		return h.hset(args)
	case hget:
		return h.hget(args)
	case hgetall:
		return h.hgetAll(args)
	case del:
		return h.del(args)
	case incr:
		return h.incrOrDecr(args, struct{ doincrement bool }{true})
	case decr:
		return h.incrOrDecr(args, struct{ doincrement bool }{false})
	case rappend:
		return h.append(args)
	case lpush:
		return h.lpush(args)

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

func (h *Handler) hset(args []resp.Value) resp.Value {
	if len(args) < 3 {
		logrus.Error("atleast one key name and key-value pair is required")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "atleast one key name and key-value pair is required"}
	}

	key := args[0]
	val := map[string]string{}

	for i := 1; i < len(args); i += 2 {
		// Ensure we don't go out of bounds
		if i+1 < len(args) {
			logrus.Debugf("Setting key %s to value %s", args[i].Bulk, args[i+1].Bulk)
			val[args[i].Bulk] = args[i+1].Bulk
		}
	}

	h.mu.Lock()
	h.hsets[key.Bulk] = val
	h.mu.Unlock()

	return okResponse()
}

func (h *Handler) hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		logrus.Error("invalid number of arguments")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "Invalid number of arguments for HGET"}
	}
	key := args[0]
	field := args[1]

	h.mu.Lock()
	val := h.hsets[key.Bulk][field.Bulk]
	h.mu.Unlock()
	return resp.Value{Typ: objects.BULK_STRING, Bulk: val}
}

func (h *Handler) hgetAll(args []resp.Value) resp.Value {
	if len(args) != 1 {
		logrus.Error("invalid number of arguments , expected only key")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "Invalid number of arguments for HGETALL ,expected only key"}
	}

	h.mu.Lock()
	val := h.hsets[args[0].Bulk]
	h.mu.Unlock()

	res := resp.Value{
		Typ: objects.ARRAY,
		Arr: make([]resp.Value, 0),
	}
	if val == nil {
		res.Arr = []resp.Value{}
		return res
	}

	arrSize := len(val) * 2
	res.Arr = make([]resp.Value, 0, arrSize)

	for k, v := range val {
		res.Arr = append(res.Arr, resp.Value{Typ: objects.BULK_STRING, Bulk: k})
		res.Arr = append(res.Arr, resp.Value{Typ: objects.BULK_STRING, Bulk: v})
	}
	return res
}

func (h *Handler) del(args []resp.Value) resp.Value {
	if len(args) == 0 {
		logrus.Error("args cannot be empty")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "args cannot be empty"}

	}
	h.mu.Lock()
	delete(h.sets, args[0].Bulk)
	return okResponse()
}

func (h *Handler) incrOrDecr(args []resp.Value, params struct{ doincrement bool }) resp.Value {
	if len(args) != 1 {
		logrus.Error("invalid args to increment")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "Invalid number of arguments,only keyname expected"}
	}

	key := args[0]

	h.mu.Lock()
	val, ok := h.sets[key.Bulk]
	if !ok {
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: ""}
	}

	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logrus.Error("error in parsing int ", err)
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: ""}
	}
	if params.doincrement {
		valInt++
	} else {
		valInt--
	}
	h.sets[key.Bulk] = strconv.FormatInt(valInt, 10)
	h.mu.Unlock()
	return resp.Value{Typ: objects.INTEGER, Num: int(valInt)}
}

func (h *Handler) append(args []resp.Value) resp.Value {
	if len(args) != 2 {
		logrus.Error("key and value to be appended to the string is needed")
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: "key and value to be appended to the string is needed"}
	}

	key := args[0]
	val := args[1]

	h.mu.Lock()

	valStr, ok := h.sets[key.Bulk]
	if !ok {
		h.mu.Unlock()
		return resp.Value{Typ: objects.SIMPLE_STRING, Str: ""}
	}

	valStr += val.Bulk

	h.sets[key.Bulk] = valStr
	h.mu.Unlock()
	return resp.Value{Typ: objects.SIMPLE_STRING, Str: valStr}
}

func (h *Handler) lpush(args []resp.Value) resp.Value {
	key := args[0]
	h.mu.Lock()

	val, ok := h.lists[key.Bulk]
	if !ok {
		logrus.Info("key doesn't exist")
		val = make([]string, 0)
	}

	l := []string{}
	l = append(l, val...)

	for i := 1; i < len(args); i++ {
		utils.Insert(l, 0, args[i].Bulk)
	}
	h.lists[key.Bulk] = l
	h.mu.Unlock()
	return resp.Value{Typ: objects.INTEGER, Num: len(l)}
}
