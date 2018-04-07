package server

import "flag"

//Config represents server configuration
type Config interface {
	Load() error
	Port() string
	TLSPair() (string, string)
}

//DefConfig represents the structure of the server
type DefConfig struct {
	//Server Config Settings
	port    string
	tlsCert string
	tlsKey  string
}

//NewConfig creates a default config file
func NewConfig() *DefConfig {
	return &DefConfig{
		port: "8080",
	}
}

//Load gets config from flags
func (c *DefConfig) Load() error {
	flag.StringVar(&c.port, "port", "8080", "Server port, default 8080")
	flag.StringVar(&c.tlsCert, "tlscert", "", "Path to TLS Certificate")
	flag.StringVar(&c.tlsKey, "tlskey", "", "Path to TLS Key")
	flag.Parse()

	return nil
}

//TLSPair returns the tls cert and key
func (c *DefConfig) TLSPair() (string, string) {
	return c.tlsCert, c.tlsKey
}

//Port returns the port to run the server on
func (c *DefConfig) Port() string {
	return c.port
}
