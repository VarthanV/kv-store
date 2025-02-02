package command

import (
	"github.com/VarthanV/kv-store/pkg/objects"
	"github.com/VarthanV/kv-store/resp"
)

type Handler map[string]func(args []resp.Value) resp.Value

func ping(args []resp.Value) resp.Value {
	return resp.Value{Typ: objects.SIMPLE_STRING, Str: "PONG"}
}

var Handlers = Handler{
	"ping": ping,
}
