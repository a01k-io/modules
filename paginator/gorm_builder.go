package paginator

import (
	"strings"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func mapSortByToDefault(sortBy []string) []clause.OrderByColumn {
	columns := make([]clause.OrderByColumn, 0)
	for _, sort := range sortBy {
		parts := strings.Split(sort, ":")
		fieldName := parts[0]
		desc := false
		if len(parts) > 1 && strings.ToLower(parts[1]) == "desc" {
			desc = true
		}

		columnName := schema.NamingStrategy{}.ColumnName("", fieldName)
		columns = append(columns, clause.OrderByColumn{Column: clause.Column{Name: columnName}, Desc: desc})
	}
	return columns
}

// BuildPaginationQuery builds pagination SQL clauses without directly using gorm.DB.
func BuildPaginationQuery(params PaginationQueryParam) []clause.Expression {
	sortByColumns := mapSortByToDefault(params.SortBy)
	clauses := make([]clause.Expression, 0)
	if len(sortByColumns) > 0 {
		clauses = append(clauses, clause.OrderBy{
			Columns: sortByColumns,
		})
	}
	if params.PageSize > 0 {
		offset := 0
		if params.PageNo > 1 {
			offset = int((params.PageNo - 1) * params.PageSize)
		}

		clauses = append(clauses, clause.Limit{
			Limit:  params.PageSizePInt(),
			Offset: offset,
		})
	}

	return clauses
}
