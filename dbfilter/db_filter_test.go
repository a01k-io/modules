package dbfilter

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/a01k-io/modules/paginator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.openly.dev/pointy"
)

func TestQueryBuilderToMongo(t *testing.T) {

	cases := []struct {
		name    string
		qb      QueryBuilder
		opts    *options.FindOptions
		want    bson.M
		wantErr bool
	}{
		{
			name: "query builder query without options",
			qb: QueryBuilder{
				Query: bson.M{
					"tenant_id":   "tenant1",
					"facility_id": "facility1",
				},
			},
			opts: &options.FindOptions{},
			want: bson.M{
				"tenant_id":   "tenant1",
				"facility_id": "facility1",
			},
			wantErr: false,
		},
		{
			name: "query builder query with all options",
			qb: QueryBuilder{
				Query: bson.M{
					"tenant_id":   "tenant1",
					"facility_id": "facility1",
				},
				Sort: []SortType{
					{Name: "_id", Direction: Desc},
					{Name: "tenant_id", Direction: Asc},
				},
				Limit: 1,
				Skip:  10,
			},
			opts: &options.FindOptions{
				Limit: pointy.Int64(1),
				Skip:  pointy.Int64(10),
				Sort: bson.D{
					bson.E{Key: "_id", Value: Desc},
					bson.E{Key: "tenant_id", Value: Asc},
				},
			},
			want: bson.M{
				"tenant_id":   "tenant1",
				"facility_id": "facility1",
			},
			wantErr: false,
		},
		{
			name: "pagination query converted to mongo",
			qb: QueryBuilder{
				Query: bson.M{
					"tenant_id": "tenant1",
				},
				Sort: []SortType{
					{Name: "tenant_id", Direction: Desc},
				},
				Limit: 20,
				Skip:  20,
			},
			opts: &options.FindOptions{
				Limit: pointy.Int64(20),
				Skip:  pointy.Int64(20),
				Sort: bson.D{
					bson.E{Key: "tenant_id", Value: Desc},
				},
			},
			want: bson.M{
				"tenant_id": "tenant1",
			},
		},
		{
			name: "error not supported query type",
			qb: QueryBuilder{
				Query: bson.D{},
			},
			opts:    &options.FindOptions{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotOpts options.FindOptions
			got, err := c.qb.ToMongo(&gotOpts)
			if (err != nil) != c.wantErr {
				t.Errorf("QueryBuilder.ToMongo() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			assert.Equal(t, c.opts, &gotOpts)

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("QueryBuilder.ToMongo() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestGetDBFilterPaginationQuery(t *testing.T) {
	cases := []struct {
		name                string
		filter              bson.M
		paginationParameter paginator.PaginationQueryParam
		want                *QueryBuilder
		wantErr             *errorops.Error
	}{
		{
			name: "pagination query without any other parameter other than page size",
			paginationParameter: paginator.PaginationQueryParam{
				PageSize: 20,
			},
			filter: bson.M{
				"tenant_id": "tenant1",
			},
			want: &QueryBuilder{
				Query: bson.M{
					"tenant_id": "tenant1",
				},
				Sort:  []SortType{{Name: "_id", Direction: Desc}},
				Skip:  0,
				Limit: 20,
			},
		},
		{
			name: "pagination query with all pagination parameters, skip ignored when id is provided",
			filter: bson.M{
				"tenant_id": "tenant1",
			},
			paginationParameter: paginator.PaginationQueryParam{
				PageNo:   2,
				PageSize: 20,
				LastID:   "61f126a1cf897aa26118d344",
				Type:     paginator.NextPage,
				SortBy:   []string{"tenant_id:asc"},
			},
			want: &QueryBuilder{
				Query: bson.M{
					"tenant_id": "tenant1",
					"_id": bson.M{
						"$lt": primitive.ObjectID{0x61, 0xf1, 0x26, 0xa1, 0xcf, 0x89, 0x7a, 0xa2, 0x61, 0x18, 0xd3, 0x44},
					},
				},
				Sort:  []SortType{{Name: "tenant_id", Direction: Asc}},
				Limit: 20,
				Skip:  0,
			},
		},
		{
			name: "pagination query with all pagination parameters, no last id given",
			filter: bson.M{
				"tenant_id": "tenant1",
			},
			paginationParameter: paginator.PaginationQueryParam{
				PageNo:   2,
				PageSize: 20,
				Type:     paginator.NextPage,
				SortBy:   []string{"tenant_id:desc"},
			},
			want: &QueryBuilder{
				Query: bson.M{
					"tenant_id": "tenant1",
				},
				Sort:  []SortType{{Name: "tenant_id", Direction: Desc}},
				Limit: 20,
				Skip:  20,
			},
		},
		{
			name: "pagination query with all pagination parameters, invalid sort format",
			filter: bson.M{
				"tenant_id": "tenant1",
			},
			paginationParameter: paginator.PaginationQueryParam{
				PageNo:   2,
				PageSize: 20,
				Type:     paginator.NextPage,
				SortBy:   []string{"invalid_sort_format"},
			},
			wantErr: &errorops.Error{
				Message: "bad request! fields validation error exist",
				Fields: []errorops.Field{
					{
						Name:    "sort_by",
						Message: []string{"invalid sort_by value: invalid_sort_format"},
						Code:    "code-invalid",
					},
				},
				Code:       "invalid-fields",
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid pagination parameter last id",
			paginationParameter: paginator.PaginationQueryParam{
				LastID: "61f126a1cf897aa26118d34461f126a1cf897aa26118d344",
			},
			wantErr: &errorops.Error{
				Message: "unable to parse last_id",
				Code:    "processing",
				Err:     primitive.ErrInvalidHex,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, gotErr := BuildQuery(c.filter, c.paginationParameter)
			if gotErr != nil || c.wantErr != nil {
				c.wantErr.ID = gotErr.ID
				assert.Equal(t, c.wantErr, gotErr)
			} else {
				assert.Equal(t, c.want, got)
			}
		})
	}
}
