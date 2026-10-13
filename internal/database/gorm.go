package database

import "gorm.io/gorm/clause"

func toClauseColumns(columns []string) []clause.Column {
	result := make([]clause.Column, len(columns))

	for i, column := range columns {
		result[i] = clause.Column{
			Name: column,
		}
	}

	return result
}
