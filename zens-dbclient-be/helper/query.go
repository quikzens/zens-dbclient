package helper

import "zens-db/entity"

func BuildOrderQuery(param entity.OrderQueryParam) string {
	orderQuery := ""
	if param.SortBy != "" {
		orderQuery = param.SortBy + " DESC"
		if param.OrderBy == "asc" || param.OrderBy == "ASC" {
			orderQuery = param.SortBy + " ASC"
		}
	}
	return orderQuery
}
