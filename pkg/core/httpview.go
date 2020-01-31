package core

import (
	"net/http"
)

type HTTPView interface {
	Process(pisigContext *PisigContext) http.HandlerFunc
}
