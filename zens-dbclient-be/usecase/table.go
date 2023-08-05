package usecase

import (
	"context"
	"zens-db/entity"
)

func (u *Usecase) GetTables(ctx context.Context, connectionId int) ([]string, error) {
	tableNames, err := u.repo.GetTables(connectionId)
	if err != nil {
		return []string{}, err
	}

	return tableNames, nil
}

func (u *Usecase) GetTableColumns(ctx context.Context, connectionId int, tableName string) ([]map[string]any, error) {
	tableNames, err := u.repo.GetTableColumns(connectionId, tableName)
	if err != nil {
		return []map[string]any{}, err
	}

	return tableNames, nil
}

func (u *Usecase) GetTableRecords(ctx context.Context, connectionId int, param entity.GetTableRecordsParam) (entity.GetTableRecordsResult, error) {
	tableRecords, count, err := u.repo.GetTableRecords(connectionId, param)
	if err != nil {
		return entity.GetTableRecordsResult{}, err
	}

	return entity.GetTableRecordsResult{
		Data:   tableRecords,
		Total:  count,
		Limit:  param.Limit,
		Offset: param.Offset,
	}, nil
}
