package xenserver

import (
	xenapi "github.com/terra-farm/go-xen-api-client"
)

// Config ...
type Config struct {
	URL      string
	Username string
	Password string
}

// Connection ...
type Connection struct {
	client  *xenapi.Client
	session xenapi.SessionRef
}

// NewConnection ...
func (cfg *Config) NewConnection() (*Connection, error) {
	client, err := xenapi.NewClient(cfg.URL, nil)
	if err != nil {
		return nil, err
	}

	session, err := client.Session.LoginWithPassword(cfg.Username, cfg.Password, "1.0", "terraform")
	if err != nil {
		return nil, err
	}

	return &Connection{client, session}, nil
}
