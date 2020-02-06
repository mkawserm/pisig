package view

import (
	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/golang/glog"
	"github.com/mkawserm/pisig/pkg/core"
	"net/http"
	"net/textproto"
	"reflect"
	"unsafe"
)

const (
	toLower = 'a' - 'A' // for use with OR.
	//toUpper  = ^byte(toLower) // for use with AND.
	toLower8 = uint64(toLower) |
		uint64(toLower)<<8 |
		uint64(toLower)<<16 |
		uint64(toLower)<<24 |
		uint64(toLower)<<32 |
		uint64(toLower)<<40 |
		uint64(toLower)<<48 |
		uint64(toLower)<<56
)

func StrToBytes(str string) (bts []byte) {
	s := (*reflect.StringHeader)(unsafe.Pointer(&str))
	b := (*reflect.SliceHeader)(unsafe.Pointer(&bts))
	b.Data = s.Data
	b.Len = s.Len
	b.Cap = s.Len
	return
}

func StrHasToken(header, token string) (has bool) {
	return BtsHasToken(StrToBytes(header), StrToBytes(token))
}

func BtsHasToken(header, token []byte) (has bool) {
	httphead.ScanTokens(header, func(v []byte) bool {
		has = BtsEqualFold(v, token)
		return !has
	})
	return
}

// StrEqualFold checks s to be case insensitive equal to p.
// Note that p must be only ascii letters. That is, every byte in p belongs to
// range ['a','z'] or ['A','Z'].
func StrEqualFold(s, p string) bool {
	return BtsEqualFold(StrToBytes(s), StrToBytes(p))
}

// btsEqualFold checks s to be case insensitive equal to p.
// Note that p must be only ascii letters. That is, every byte in p belongs to
// range ['a','z'] or ['A','Z'].
func BtsEqualFold(s, p []byte) bool {
	if len(s) != len(p) {
		return false
	}
	n := len(s)
	// Prepare manual conversion on bytes that not lay in uint64.
	m := n % 8
	for i := 0; i < m; i++ {
		if s[i]|toLower != p[i]|toLower {
			return false
		}
	}
	// Iterate over uint64 parts of s.
	n = (n - m) >> 3
	if n == 0 {
		// There are no more bytes to compare.
		return true
	}

	for i := 0; i < n; i++ {
		x := m + (i << 3)
		av := *(*uint64)(unsafe.Pointer(&s[x]))
		bv := *(*uint64)(unsafe.Pointer(&p[x]))
		if av|toLower8 != bv|toLower8 {
			return false
		}
	}

	return true
}

var (
	//headerHost          = "Host"
	headerUpgrade       = "Upgrade"
	headerConnection    = "Connection"
	headerAuthorization = "Authorization"
	//headerSecVersion  = "Sec-WebSocket-Version"
	headerSecProtocol = "Sec-WebSocket-Protocol"
	//headerSecExtensions = "Sec-WebSocket-Extensions"
	//headerSecKey        = "Sec-WebSocket-Key"
	//headerSecAccept     = "Sec-WebSocket-Accept"

	//headerHostCanonical          = textproto.CanonicalMIMEHeaderKey(headerHost)
	headerUpgradeCanonical    = textproto.CanonicalMIMEHeaderKey(headerUpgrade)
	headerConnectionCanonical = textproto.CanonicalMIMEHeaderKey(headerConnection)
	//headerAuthorizationCanonical = textproto.CanonicalMIMEHeaderKey(headerAuthorization)
	//headerSecVersionCanonical    = textproto.CanonicalMIMEHeaderKey(headerSecVersion)
	//headerSecProtocolCanonical = textproto.CanonicalMIMEHeaderKey(headerSecProtocol)
	//headerSecExtensionsCanonical = textproto.CanonicalMIMEHeaderKey(headerSecExtensions)
	//headerSecKeyCanonical        = textproto.CanonicalMIMEHeaderKey(headerSecKey)
	//headerSecAcceptCanonical     = textproto.CanonicalMIMEHeaderKey(headerSecAccept)
)

//const (
//	SuperGroup   = "super"
//	ServiceGroup = "service"
//	NormalGroup  = "normal"
//)
//
//const WebSocketConnection = 1
//const WebSocketSecureConnection = 2
//
//const HTTPConnection = 3
//const HTTPSConnection = 4
//
//const (
//	WebSocketOpContinuation byte = 0x0
//	WebSocketOpText         byte = 0x1
//	WebSocketOpBinary       byte = 0x2
//	WebSocketOpClose        byte = 0x8
//	WebSocketOpPing         byte = 0x9
//	WebSocketOpPong         byte = 0xa
//)

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
		defer func() {
			if glog.V(3) {
				glog.Infof("WebSocketView - Process - http.HandlerFunc done\n")
			}
		}()

		if glog.V(3) {
			glog.Infof("WebSocketView - Process - http.HandlerFunc begin\n")
			glog.Infof("PATH - " + request.URL.Path + "\n")
			glog.Infof("Checking CORS")
		}

		if !v.pisig.PisigContext().CORSOptions.CROSCheckAllowNext(writer, request) {
			if glog.V(1) {
				glog.Infof("CORS block!!!\n")
			}
			return
		}

		if glog.V(3) {
			glog.Infof("Processing middleware\n")
		}

		for _, httpMiddlewareView := range v.pisig.MiddlewareViewList() {
			next := httpMiddlewareView.ProcessAllowNext(v.pisig)(writer, request)
			if next == false {
				return
			}
		}

		if glog.V(3) {
			glog.Infof("Middleware processed\n")
		}
		// BOILERPLATE END

		// MAIN LOGIC
		if glog.V(3) {
			glog.Infof("Executing websocket upgrade logic\n")
		}

		if request.Method != http.MethodGet {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			return
		}

		if request.ProtoMajor < 1 || (request.ProtoMajor == 1 && request.ProtoMinor < 1) {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			return
		}

		if request.Host == "" {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			return
		}

		if u := httpGetHeader(request.Header, headerUpgradeCanonical); u != "websocket" &&
			!StrEqualFold(u, "websocket") {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			return
		}

		if c := httpGetHeader(request.Header, headerConnectionCanonical); c != "Upgrade" &&
			!StrHasToken(c, "upgrade") {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			return
		}

		conn, _, _, err := ws.UpgradeHTTP(request, writer)

		if err != nil {
			glog.Errorf("Failed to do upgrade handshake - %v\n", err)
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			return
		}

		if !v.pisig.AddWebSocketConnection(conn) {
			glog.Errorf("Failed to add connection\n")
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(v.pisig.PisigContext().PisigMessage.HTTP400())
			_ = conn.Close()
			return
		}

		if glog.V(3) {
			glog.Infof("Websocket upgrade logic executed\n")
		}
		// MAIN LOGIC END

	}
}

func httpGetHeader(h http.Header, key string) string {
	if h == nil {
		return ""
	}
	v := h[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
