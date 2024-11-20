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

		// Convert fieldName to GORM default naming convention (snake_case)
		columnName := schema.NamingStrategy{}.ColumnName("", fieldName)
		columns = append(columns, clause.OrderByColumn{Column: clause.Column{Name: columnName}, Desc: desc})
	}
	return columns
}

// BuildPaginationQuery builds pagination SQL clauses without directly using gorm.DB.
func BuildPaginationQuery(params PaginationQueryParam) []clause.Expression {
	// Function to map SortBy to database columns using GORM naming conventions
	mapSortByToDefault := func(sortBy []string) []clause.OrderByColumn {
		columns := []clause.OrderByColumn{}
		for _, sort := range sortBy {
			parts := strings.Split(sort, ":")
			fieldName := parts[0]
			desc := false
			if len(parts) > 1 && strings.ToLower(parts[1]) == "desc" {
				desc = true
			}

			// Convert fieldName to GORM default naming convention (snake_case)
			columnName := schema.NamingStrategy{}.ColumnName("", fieldName)
			columns = append(columns, clause.OrderByColumn{Column: clause.Column{Name: columnName}, Desc: desc})
		}
		return columns
	}

	// Map SortBy to clause.OrderByColumn
	sortByColumns := mapSortByToDefault(params.SortBy)

	// Build clauses
	clauses := make([]clause.Expression, 0)

	// Add sorting
	if len(sortByColumns) > 0 {
		clauses = append(clauses, clause.OrderBy{
			Columns: sortByColumns,
		})
	}

	// Add limit and offset
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
