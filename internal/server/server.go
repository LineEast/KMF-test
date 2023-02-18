package server

import (
	"kmfRedirect/internal/database"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type (
	Server struct {
		router   *router.Router
		database *database.Database

		configuration *Configuration
	}

	Configuration struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	}
)

func New(database *database.Database, configuration *Configuration) *Server {
	server := &Server{
		database:      database,
		configuration: configuration,
	}

	r := router.New()
	r.POST("/", server.Main())

	server.router = r

	return server
}

func (s *Server) Run() error {
	d := fasthttp.Server{
		Handler:           s.router.Handler,
		StreamRequestBody: true,
	}

	return d.ListenAndServe(s.configuration.Host + ":" + s.configuration.Port)
}
