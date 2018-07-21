package h2conn

import (
	"fmt"
	"io"
	"net/http"
)

// ErrHTTP2NotSupported is returned by Accept if the client connection does not
// support HTTP2 connection.
// The server than can response to the client with an HTTP1.1 as he wishes.
var ErrHTTP2NotSupported = fmt.Errorf("HTTP2 not supported")

// Server can "accept" an http2 connection to obtain a read/write object
// for full duplex communication with a client.
type Server struct {
	StatusCode int
}

var defaultUpgrader = Server{
	StatusCode: http.StatusOK,
}

func Accept(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	return defaultUpgrader.Accept(w, r)
}

// Accept is used on a server http.Handler.
// It handles a request and "upgrade" the request connection to a websocket-like
// full-duplex communication.
// If the client does not support HTTP2, an ErrHTTP2NotSupported is returned.
//
// Usage:
//
//      func (w http.ResponseWriter, r *http.Request) {
//          conn, err := h2conn.Accept(w, r)
//          if err != nil {
//		        log.Printf("Failed creating http2 connection: %s", err)
//		        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		        return
//	        }
//          // use conn
//      }
func (u *Server) Accept(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	if !r.ProtoAtLeast(2, 0) {
		return nil, ErrHTTP2NotSupported
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrHTTP2NotSupported
	}

	c := newConn(r.Context(), r.Body, &flushWrite{w: w, f: flusher})

	w.WriteHeader(u.StatusCode)
	flusher.Flush()

	return c, nil
}

type flushWrite struct {
	w io.Writer
	f http.Flusher
}

func (w *flushWrite) Write(data []byte) (int, error) {
	n, err := w.w.Write(data)
	w.f.Flush()
	return n, err
}

func (w *flushWrite) Close() error {
	return nil
}
