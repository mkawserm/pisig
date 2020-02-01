package core

import (
	"net/http"
)

type HTTPView interface {
	Process(pisig *Pisig) http.HandlerFunc
}
