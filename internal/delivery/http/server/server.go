package server

import (
	"net/http"
	"time"

	"github.com/andrsj/go-rabbit-image/internal/delivery/http/handler"
)

// New function is the constructor for HTTP Server
//
// # This server was used for graceful shutdown
//
// I actually don't know how to graceful shutdown the gin.Engine,
// so I directly shutdown the http.Server.
func New(h *handler.Handler) *http.Server {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           h.GetGinEngine(),
		ReadHeaderTimeout: time.Second,
	}

	return server
}
