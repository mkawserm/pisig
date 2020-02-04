package view

import (
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/core"
	"net/http"
)

type WebSocketView struct {
	pisig *core.Pisig
}

func (v *WebSocketView) Process(pisig *core.Pisig) http.HandlerFunc {

	// BOILERPLATE BEGIN
	if glog.V(3) {
		glog.Infof("WebSocketView - Process begin\n")
	}

	v.pisig = pisig

	if glog.V(3) {
		glog.Infof("WebSocketView - Process done\n")
	}
	// BOILERPLATE END

	return func(writer http.ResponseWriter, request *http.Request) {

		// BOILERPLATE BEGIN
		if glog.V(3) {
			glog.Infof("WebSocketView - Process - http.HandlerFunc begin\n")
			glog.Infof("PATH - " + request.URL.Path + "\n")
			glog.Infof("Checking CORS")
		}

		if !v.pisig.CORSOptions().CROSCheckAllowNext(writer, request) {
			if glog.V(1) {
				glog.Infof("CORS block!!!\n")
			}

			if glog.V(3) {
				glog.Infof("WebSocketView - Process - http.HandlerFunc done\n")
			}
			return
		}

		if glog.V(3) {
			glog.Infof("Processing middleware\n")
		}

		for _, httpMiddlewareView := range v.pisig.MiddlewareViewList() {
			next := httpMiddlewareView.ProcessAllowNext(v.pisig)(writer, request)
			if next == false {
				if glog.V(3) {
					glog.Infof("WebSocketView - Process - http.HandlerFunc done\n")
				}
				return
			}
		}

		if glog.V(3) {
			glog.Infof("WebSocketView processing complete\n")
		}
		// BOILERPLATE END

		// MAIN LOGIC
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		_, _ = writer.Write(v.pisig.PisigMessage().HTTP404())
		// MAIN LOGIC END

		// BOILERPLATE BEGIN
		if glog.V(3) {
			glog.Infof("WebSocketView - Process - http.HandlerFunc done\n")
		}
		// BOILERPLATE END
	}
}
