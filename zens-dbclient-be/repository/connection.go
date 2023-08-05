package repository

import (
	"context"
	"zens-db/entity"
)

func (r *Repository) GetConnections(ctx context.Context) []entity.Connection {
	return r.connections
}
