package core

import (
	"net/http"
)

type MiddlewareHandlerFunc func(http.ResponseWriter, *http.Request) bool

type HTTPMiddlewareView interface {
	ProcessAllowNext(pisig *Pisig) MiddlewareHandlerFunc
}
