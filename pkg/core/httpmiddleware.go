package core

import (
	"net/http"
)

type HTTPMiddleware interface {
	ProcessAllowNext(pisigContext *PisigContext, w http.ResponseWriter, r *http.Request) bool
}
