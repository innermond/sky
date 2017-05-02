package http

import (
	"net"
	"net/http"
)

type IndexServer struct {
	ln      net.Listener
	Addr    string
	Handler http.Handler
}

func (s *IndexServer) Open() error {
	ln, err := net.Listen("tcp", s.Addr)
	s.ln = ln
	if err != nil {
		return err
	}
	go func() { http.Serve(s.ln, s.Handler) }()
	return nil
}

func (s *IndexServer) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}
