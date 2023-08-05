package repository

import (
	"fmt"
	"zens-db/entity"
)

func (r *Repository) GetTables(connectionId int) ([]string, error) {
	db, err := r.getConnection(connectionId)
	if err != nil {
		return []string{}, err
	}
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name ASC"

	var tableNames = make([]string, 0)
	err = db.Raw(query).Scan(&tableNames).Error
	if err != nil {
		return []string{}, err
	}

	return tableNames, nil
}

func (r *Repository) GetTableColumns(connectionId int, tableName string) ([]map[string]any, error) {
	db, err := r.getConnection(connectionId)
	if err != nil {
		return []map[string]any{}, err
	}
	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = '%v' ORDER BY ordinal_position", tableName)

	var columns = make([]map[string]any, 0)
	err = db.Raw(query).Scan(&columns).Error
	if err != nil {
		return []map[string]any{}, err
	}

	return columns, nil
}

func (r *Repository) GetTableRecords(connectionId int, param entity.GetTableRecordsParam) ([]map[string]any, int, error) {
	db, err := r.getConnection(connectionId)
	if err != nil {
		return []map[string]any{}, 0, err
	}
	query := fmt.Sprintf("SELECT * FROM %v", param.TableName)
	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM %v", param.TableName)
	if len(param.Conditions) > 0 {
		whereQuery := buildWhereQuery(param.Conditions)
		query += whereQuery
		queryCount += whereQuery
	}
	if param.SortBy != "" {
		query += buildOrderQuery(entity.OrderQueryParam{
			SortBy:  param.SortBy,
			OrderBy: param.OrderBy,
		})
	}
	if param.Limit > 0 {
		limitQuery := buildLimitQuery(param.Limit)
		query += limitQuery
		queryCount += limitQuery
	}
	if param.Offset > 0 {
		offsetQuery := buildOffsetQuery(param.Offset)
		query += offsetQuery
		queryCount += offsetQuery
	}

	var records = make([]map[string]any, 0)
	err = db.Raw(query).Scan(&records).Error
	if err != nil {
		return []map[string]any{}, 0, err
	}

	var count int64
	err = db.Debug().Raw(queryCount).Count(&count).Error
	if err != nil {
		return []map[string]any{}, 0, err
	}

	return records, int(count), nil
}
