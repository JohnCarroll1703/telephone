package server

import (
	"net/http"
	"telephone/internal/config"
)

type Server struct {
	server *http.Server
}

func NewServer(
	config *config.Config,
) {

}
