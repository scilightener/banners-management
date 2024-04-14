package config

import "fmt"

// HTTPServer contains the settings for the HTTP server.
type HTTPServer struct {
	Address     string   `json:"address"`
	Timeout     Duration `json:"timeout"`
	IdleTimeout Duration `json:"idle_timeout"`
}

func (s HTTPServer) String() string {
	return fmt.Sprintf("{Address: %s, Timeout: %v, IdleTimeout: %v}", s.Address, s.Timeout, s.IdleTimeout)
}
