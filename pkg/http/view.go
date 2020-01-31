package http

import (
	"github.com/mkawserm/pisig/pkg/core"
	"net/http"
)

type View interface {
	Process(pisigContext *core.PisigContext) http.HandlerFunc
}
