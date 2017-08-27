package xenserver

import (
	"github.com/mborodin/go-xen-api-client"
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
	client, err := xenAPI.NewClient(cfg.URL, nil)
	if err != nil {
		return nil, err
	}

	session, err := client.Session.LoginWithPassword(cfg.Username, cfg.Password, "1.0", "terraform")
	if err != nil {
		return nil, err
	}

	return &Connection{client, session}, nil
}
