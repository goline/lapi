package lapi

import (
	"crypto/tls"
	"time"
)

// Config contains configuration
type Config interface {
	Bag
}

// ServerConfig stores configuration of server
type ServerConfig interface {
	// Address returns a string of TCP address to listen on, ":http" if empty
	Address() string

	// ReadTimeout returns a duration of http.Server.ReadTimeout
	ReadTimeout() time.Duration

	// WriteTimeout returns a duration of http.Server.WriteTimeout
	WriteTimeout() time.Duration

	// MaxHeaderBytes returns a duration of http.Server.MaxHeaderBytes
	MaxHeaderBytes() int

	// ReadHeaderTimeout returns a duration of http.Server.ReadHeaderTimeout
	ReadHeaderTimeout() time.Duration

	// IdleTimeout returns a duration of http.Server.IdleTimeout
	IdleTimeout() time.Duration
}

// SecureServerConfig stores configuration of HTTPS server
type SecureServerConfig interface {
	// TLSConfig returns config of http.Server.TLSConfig
	TLSConfig() *tls.Config
}
