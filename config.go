package main

import (
	"github.com/amfranz/go-xen-api-client"
)

// Config ...
type Config struct {
	URL      string
	Username string
	Password string
}

// Connection ...
type Connection struct {
	client  *xenAPI.Client
	session xenAPI.SessionRef
}

// NewConnection ...
func (cfg *Config) NewConnection() (*Connection, error) {
	client := xenAPI.NewClient(cfg.URL)

	session, err := client.Session().LoginWithPassword(cfg.Username, cfg.Password, "1.0", "terraform")
	if err != nil {
		return nil, err
	}

	return &Connection{client, session}, nil
}
