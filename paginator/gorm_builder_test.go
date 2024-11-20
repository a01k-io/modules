package paginator_test

import (
	"testing"

	"github.com/a01k-io/modules/paginator"
	"github.com/stretchr/testify/assert"
	"go.openly.dev/pointy"
	"gorm.io/gorm/clause"
)

func TestBuildPaginationQuery(t *testing.T) {
	tests := []struct {
		name            string
		params          paginator.PaginationQueryParam
		sortByColumns   []clause.OrderByColumn
		expectedClauses []clause.Expression
	}{
		{
			name: "Next page with sort and limit",
			params: paginator.PaginationQueryParam{
				PageSize: 10,
				LastID:   "5",
				Type:     paginator.NextPage,
				SortBy:   []string{"TenantID:desc", "Name:asc"},
			},
			sortByColumns: []clause.OrderByColumn{
				{Column: clause.Column{Name: "tenant_id"}, Desc: true},
				{Column: clause.Column{Name: "name"}, Desc: false},
			},
			expectedClauses: []clause.Expression{
				clause.OrderBy{
					Columns: []clause.OrderByColumn{
						{Column: clause.Column{Name: "tenant_id"}, Desc: true},
						{Column: clause.Column{Name: "name"}, Desc: false},
					},
				},
				clause.Limit{Limit: pointy.Pointer(10)},
				clause.Where{
					Exprs: []clause.Expression{
						clause.Expr{
							SQL:  "id > ?",
							Vars: []interface{}{"5"},
						},
					},
				},
			},
		},
		{
			name: "Previous page without LastID",
			params: paginator.PaginationQueryParam{
				PageSize: 5,
				Type:     paginator.PrevPage,
			},
			sortByColumns: []clause.OrderByColumn{
				{Column: clause.Column{Name: "name"}, Desc: true},
			},
			expectedClauses: []clause.Expression{
				clause.OrderBy{
					Columns: []clause.OrderByColumn{
						{Column: clause.Column{Name: "name"}, Desc: true},
					},
				},
				clause.Limit{Limit: pointy.Pointer(5)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualClauses := paginator.BuildPaginationQuery(tt.params, tt.sortByColumns)

			// Check that the number of clauses matches
			assert.Equal(t, len(tt.expectedClauses), len(actualClauses), "Mismatch in number of clauses")

			// Compare each clause
			for i := range tt.expectedClauses {
				assert.Equal(t, tt.expectedClauses[i], actualClauses[i], "Clause %d mismatch", i)
			}
		})
	}
}
