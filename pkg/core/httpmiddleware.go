package core

import (
	"net/http"
)

type MiddlewareHandlerFunc func(http.ResponseWriter, *http.Request) bool

type HTTPMiddleware interface {
	ProcessAllowNext(pisig *Pisig) MiddlewareHandlerFunc
}
