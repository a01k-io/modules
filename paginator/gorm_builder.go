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
func BuildPaginationQuery(params PaginationQueryParam, sortByColumns []clause.OrderByColumn) (clauses []clause.Expression) {
	// Add sorting
	if len(sortByColumns) > 0 {
		clauses = append(clauses, clause.OrderBy{
			Columns: sortByColumns,
		})
	}

	// Add limit
	if params.PageSize > 0 {
		clauses = append(clauses, clause.Limit{
			Limit: params.PageSizePInt(),
		})
	}

	// Add pagination logic (LastID and Type)
	if params.LastID != "" {
		filterClause := clause.Where{}
		if params.Type == NextPage {
			filterClause.Exprs = append(filterClause.Exprs, clause.Expr{
				SQL:  "id > ?",
				Vars: []interface{}{params.LastID},
			})
		} else if params.Type == PrevPage {
			filterClause.Exprs = append(filterClause.Exprs, clause.Expr{
				SQL:  "id < ?",
				Vars: []interface{}{params.LastID},
			})
		}
		clauses = append(clauses, filterClause)
	}

	return clauses
}
