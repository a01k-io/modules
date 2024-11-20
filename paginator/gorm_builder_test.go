package paginator_test

import (
	"github.com/a01k-io/modules/paginator"
	"go.openly.dev/pointy"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
)

func TestBuildPaginationQuery(t *testing.T) {
	tests := []struct {
		name            string
		params          paginator.PaginationQueryParam
		expectedClauses []clause.Expression
	}{
		{
			name: "Next page with limit and offset",
			params: paginator.PaginationQueryParam{
				PageSize: 6,
				PageNo:   2,
				SortBy:   []string{"created_at:desc"},
			},
			expectedClauses: []clause.Expression{
				clause.OrderBy{
					Columns: []clause.OrderByColumn{
						{Column: clause.Column{Name: "created_at"}, Desc: true},
					},
				},
				clause.Limit{
					Limit:  pointy.Pointer(6),
					Offset: 6, // (PageNo - 1) * PageSize
				},
			},
		},
		{
			name: "First page without offset",
			params: paginator.PaginationQueryParam{
				PageSize: 10,
				PageNo:   1,
				SortBy:   []string{"id:asc"},
			},
			expectedClauses: []clause.Expression{
				clause.OrderBy{
					Columns: []clause.OrderByColumn{
						{Column: clause.Column{Name: "id"}, Desc: false},
					},
				},
				clause.Limit{
					Limit:  pointy.Pointer(10),
					Offset: 0,
				},
			},
		},
		{
			name: "Page size only, no sorting or page number",
			params: paginator.PaginationQueryParam{
				PageSize: 5,
			},
			expectedClauses: []clause.Expression{
				clause.Limit{
					Limit:  pointy.Pointer(5),
					Offset: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualClauses := paginator.BuildPaginationQuery(tt.params)

			// Проверяем количество клауза
			assert.Equal(t, len(tt.expectedClauses), len(actualClauses), "Количество выражений не совпадает")

			// Проверяем каждое выражение
			for i := range tt.expectedClauses {
				assert.Equal(t, tt.expectedClauses[i], actualClauses[i], "Клауз %d не совпадает", i)
			}
		})
	}
}
