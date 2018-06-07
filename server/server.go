package server

import (
	"fmt"
	"log"
	"net/http"
)

//Server contains relevant data to our server
type Server struct {
	Config Config
	Router http.Handler
}

//New creates and initializes a new server
func New(r http.Handler, c Config) *Server {
	srv := &Server{Router: r, Config: c}
	srv.initializeConfig()
	return srv
}

//NewWithDefaultConfig creates and initializes a new server with the default configuration
func NewWithDefaultConfig(router http.Handler) *Server {
	srv := &Server{Router: router}
	srv.initializeDefaultConfig()
	return srv
}

//SetRouter ...
func (s *Server) SetRouter(r http.Handler) {
	s.Router = r
}

//Run starts http server
func (s *Server) Run() error {
	defer s.destroy()
	return s.runSandbox()
}

//RunTLS starts https server
func (s *Server) RunTLS() error {
	defer s.destroy()

	tlsCert, tlsKey := s.Config.TLSPair()
	if tlsCert == "" || tlsKey == "" {
		return fmt.Errorf("tls cert pair not found")
	}

	return s.runProduction()
}

func (s *Server) runSandbox() error {
	var p = s.Config.Port()
	if len(p) == 0 {
		p = "8080"
	}

	printDebugStart(p)
	return http.ListenAndServe(
		fmt.Sprintf(":%s", p),
		s.Router,
	)
}

func (s *Server) runProduction() error {
	var p = s.Config.Port()
	if len(p) == 0 {
		p = "80"
	}

	printDebugStart(p)
	var tlsCert, tlsKey = s.Config.TLSPair()
	return http.ListenAndServeTLS(
		fmt.Sprintf(":%s", p),
		tlsCert,
		tlsKey,
		s.Router,
	)
}

func (s *Server) initializeDefaultConfig() {
	s.Config = NewConfig()
	s.initializeConfig()
	return
}

func (s *Server) initializeConfig() {
	checkErr(s.Config.Load())
}

func (s *Server) destroy() {
	//nothing to destroy
}

func printDebugStart(port string) {
	log.Printf("Listening at %s...\n", port)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
