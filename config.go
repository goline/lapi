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

type FactoryConfig struct {
	Bag
	address           string
	readTimeout       time.Duration
	writeTimeout      time.Duration
	maxHeaderBytes    int
	readHeaderTimeout time.Duration
	idleTimeout       time.Duration
	tlsConfig         *tls.Config
}

func (c *FactoryConfig) Address() string {
	return c.address
}

func (c *FactoryConfig) WithAddress(address string) ServerConfig {
	c.address = address
	return c
}

func (c *FactoryConfig) ReadTimeout() time.Duration {
	return c.readTimeout
}

func (c *FactoryConfig) WithReadTimeout(timeout time.Duration) ServerConfig {
	c.readTimeout = timeout
	return c
}

func (c *FactoryConfig) WriteTimeout() time.Duration {
	return c.writeTimeout
}

func (c *FactoryConfig) WithWriteTimeout(timeout time.Duration) ServerConfig {
	c.writeTimeout = timeout
	return c
}

func (c *FactoryConfig) MaxHeaderBytes() int {
	return c.maxHeaderBytes
}

func (c *FactoryConfig) WithMaxHeaderBytes(bytes int) ServerConfig {
	c.maxHeaderBytes = bytes
	return c
}

func (c *FactoryConfig) ReadHeaderTimeout() time.Duration {
	return c.readHeaderTimeout
}

func (c *FactoryConfig) WithReadHeaderTimeout(timeout time.Duration) ServerConfig {
	c.readHeaderTimeout = timeout
	return c
}

func (c *FactoryConfig) IdleTimeout() time.Duration {
	return c.idleTimeout
}

func (c *FactoryConfig) WithIdleTimeout(timeout time.Duration) ServerConfig {
	c.idleTimeout = timeout
	return c
}

func (c *FactoryConfig) TLSConfig() *tls.Config {
	return c.tlsConfig
}

func (c *FactoryConfig) WithTLSConfig(config *tls.Config) SecureServerConfig {
	c.tlsConfig = config
	return c
}
