package server

import (
	"net/http"

	"github.com/andrsj/go-rabbit-image/internal/delivery/http/handler"
)

func New(h *handler.Handler) *http.Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: h.Engine,
	}
	return server
}
