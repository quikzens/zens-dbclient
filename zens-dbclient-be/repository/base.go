package repository

import (
	"zens-db/entity"

	"golang.org/x/exp/slices"
)

type Repository struct {
	lastIdConnection int
	connections      []entity.Connection
}

func New() *Repository {
	return &Repository{
		lastIdConnection: 0,
		connections:      []entity.Connection{},
	}
}

func (r *Repository) AddConnection(c entity.Connection) int {
	r.lastIdConnection = r.lastIdConnection + 1
	c.Id = r.lastIdConnection
	r.connections = append(r.connections, c)
	return c.Id
}

func (r *Repository) DeleteConnection(connectionId int) (int, error) {
	// find connection index
	connectionIdx := slices.IndexFunc(r.connections, func(c entity.Connection) bool {
		return c.Id == connectionId
	})

	// close connection
	conn, err := r.connections[connectionIdx].Client.DB()
	if err != nil {
		return connectionId, err
	}
	err = conn.Close()
	if err != nil {
		return connectionId, err
	}

	// remove connection from list
	r.connections = slices.Delete(r.connections, connectionIdx, connectionIdx+1)
	return connectionId, nil
}
