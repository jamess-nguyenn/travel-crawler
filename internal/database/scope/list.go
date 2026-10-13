package scope

import (
	"strings"
	"travel-crawler/internal/request"
)

func GetSortColumn(
	sortBy string,
	sortColumns map[string]string,
	defaultSortColumn string,
) string {
	if sortColumn, ok := sortColumns[sortBy]; ok {
		return sortColumn
	}

	return defaultSortColumn
}

func BuildListScopes(
	req request.SortRequest,
	sortColumns map[string]string,
	defaultSortColumn string,
) []Scope {
	scopes := make([]Scope, 0)

	sortColumn := GetSortColumn(
		req.GetSortBy(),
		sortColumns,
		defaultSortColumn,
	)

	sortBys := strings.Split(sortColumn, ",")
	for _, sortBy := range sortBys {
		scopes = append(scopes, Order(sortBy, req.GetSortOrder()))
	}

	if pagination, ok := req.(request.PaginationRequest); ok && !pagination.NoLimit() {
		scopes = append(scopes, Paginate(pagination.GetPage(), pagination.GetPageSize()))
	}

	return scopes
}
