package usecase

import (
	"context"
	"fmt"
	"zens-db/entity"
	"zens-db/helper"
)

func (u *Usecase) GetConnections(ctx context.Context) []entity.Connection {
	return u.repo.GetConnections(ctx)
}

func (u *Usecase) CreateConnection(ctx context.Context, param entity.CreateConnectionParam) (entity.CreateConnectionResult, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", param.User, param.Password, param.Host, param.Port, param.DatabaseName)
	db, err := u.initDbConnection(dsn)
	if err != nil {
		return entity.CreateConnectionResult{}, helper.HandleDbErr(err)
	}

	connection := entity.Connection{
		Client: db,
		Credential: entity.Credential{
			Host:         param.Host,
			Port:         param.Port,
			DatabaseName: param.DatabaseName,
			User:         param.User,
			Password:     param.Password,
		},
	}

	return entity.CreateConnectionResult{
		ConnectionId: u.repo.AddConnection(connection),
	}, nil
}

func (u *Usecase) DeleteConnection(ctx context.Context, connectionId int) (entity.DeleteConnectionResult, error) {
	connectionId, err := u.repo.DeleteConnection(connectionId)
	if err != nil {
		return entity.DeleteConnectionResult{}, err
	}

	return entity.DeleteConnectionResult{ConnectionId: connectionId}, nil
}
