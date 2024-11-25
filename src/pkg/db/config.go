//nolint:tagalign // Tags of config params are manually aligned.
package db

import (
	"context"
)

type DatabaseConfig struct {
	DatabaseWriteURL string `param:"database-write-url" default:""  desc:"database connection string for write operations"`
	DatabaseReadURL  string `param:"database-read-url"  default:"" desc:"database connection string for read operations"`
}

type Clients struct {
	Write Client
	Read  Client
}

func NewClients(ctx context.Context, cfg *DatabaseConfig, debug bool) *Clients {
	return &Clients{
		Write: newClient(ctx, "write", cfg.DatabaseWriteURL, debug),
		Read:  newClient(ctx, "read", cfg.DatabaseReadURL, debug),
	}
}

func CloseClients(clients *Clients) {
	closeClient("read", clients.Read)
	closeClient("write", clients.Write)
}
