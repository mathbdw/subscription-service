package httpserver

import (
	"fmt"
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Address -.
func Address(host string, port uint16) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort(host, fmt.Sprintf("%d", port))
	}
}

// Prefork -.
func Prefork(prefork bool) Option {
	return func(s *Server) {
		s.prefork = prefork
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
