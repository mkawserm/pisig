package http

import (
	"github.com/mkawserm/pisig/pkg/core"
	"net/http"
)

type Middleware interface {
	ProcessAllowNext(pisigContext *core.PisigContext, w http.ResponseWriter, r *http.Request) bool
}
