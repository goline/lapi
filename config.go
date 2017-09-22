package lapi

import (
	"crypto/tls"
	"time"
)

// ServerConfig stores configuration of server
type ServerConfig interface {
	// Address returns a string of TCP address to listen on, ":http" if empty
	Address() string

	// WithAddress sets address string
	WithAddress(address string) ServerConfig

	// ReadTimeout returns a duration of http.Server.ReadTimeout
	ReadTimeout() time.Duration

	// WithReadTimeout sets ReadTimeout
	WithReadTimeout(timeout time.Duration) ServerConfig

	// WriteTimeout returns a duration of http.Server.WriteTimeout
	WriteTimeout() time.Duration

	// WithWriteTimeout sets WriteTimeout
	WithWriteTimeout(timeout time.Duration) ServerConfig

	// MaxHeaderBytes returns a duration of http.Server.MaxHeaderBytes
	MaxHeaderBytes() int

	// WithMaxHeaderBytes sets MaxHeaderBytes
	WithMaxHeaderBytes(bytes int) ServerConfig

	// ReadHeaderTimeout returns a duration of http.Server.ReadHeaderTimeout
	ReadHeaderTimeout() time.Duration

	// WithReadHeaderTimeout sets ReadHeaderTimeout
	WithReadHeaderTimeout(timeout time.Duration) ServerConfig

	// IdleTimeout returns a duration of http.Server.IdleTimeout
	IdleTimeout() time.Duration

	// WithIdleTimeout sets IdleTimeout
	WithIdleTimeout(timeout time.Duration) ServerConfig
}

// SecureServerConfig stores configuration of HTTPS server
type SecureServerConfig interface {
	// TLSConfig returns config of http.Server.TLSConfig
	TLSConfig() *tls.Config

	// WithTLSConfig sets TLSConfig
	WithTLSConfig(config *tls.Config) SecureServerConfig
}
