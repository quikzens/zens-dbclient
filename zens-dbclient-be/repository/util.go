package repository

import (
	"fmt"
	"zens-db/entity"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func buildOrderQuery(param entity.OrderQueryParam) string {
	orderQuery := ""
	if param.SortBy != "" {
		orderQuery = param.SortBy + " DESC"
		if param.OrderBy == "asc" || param.OrderBy == "ASC" {
			orderQuery = param.SortBy + " ASC"
		}
	}
	return fmt.Sprintf(" ORDER BY %s", orderQuery)
}

func buildWhereQuery(conditions []entity.Condition) string {
	whereQuery := ""
	for i, condition := range conditions {
		if i == 0 {
			whereQuery += fmt.Sprintf(" WHERE %v %v", condition.Field, condition.Operator)
		} else {
			whereQuery += fmt.Sprintf(" AND %v %v", condition.Field, condition.Operator)
		}
		if condition.FirstValue != "" {
			whereQuery += fmt.Sprintf(" '%v'", condition.FirstValue)
		}
		if condition.SecondValue != "" { // for 'BETWEEN' operator
			whereQuery += fmt.Sprintf(" AND '%v'", condition.SecondValue)
		}
	}
	return whereQuery
}

func buildLimitQuery(limit int) string {
	return fmt.Sprintf(" LIMIT %d", limit)
}

func buildOffsetQuery(offset int) string {
	return fmt.Sprintf(" OFFSET %d", offset)
}

func (r *Repository) getConnection(connectionId int) (*gorm.DB, error) {
	connectionIndex := slices.IndexFunc(r.connections, func(c entity.Connection) bool { return c.Id == connectionId })
	if connectionIndex < 0 {
		return nil, entity.ConnectionNotFoundError{}
	}
	return r.connections[connectionIndex].Client, nil
}
