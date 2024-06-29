package dbfilter

import (
	"errors"
	"fmt"

	"github.com/a01k-io/modules/paginator"
	"github.com/a01k-io/modules/stringops"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SortDirection is a filter operator for sorting using the field.
type SortDirection int

const (
	// Asc is ascending sort order.
	Asc SortDirection = 1
	// Desc is decending sort order.
	Desc SortDirection = -1
	// None when sort is not defined.
	None SortDirection = 0
)

// SortType is a filter operator for sorting using the field.
type SortType struct {
	Name      string
	Direction SortDirection
}

type Operator string

const (
	// LTOp less than operator.
	LTOp Operator = "$lt"

	// GTOp greater than operator.
	GTOp Operator = "$gt"
)

type QueryBuilder struct {
	Query interface{}
	Sort  []SortType
	Limit int64
	Skip  int64
}

var (
	// ErrInvalidShortType means the short type is invalid.
	ErrInvalidShortType      = errors.New("invalid sort type")
	ErrorUnableToParseLastID = errors.New("unable to parse last_id")
	//PaginationTypeQueryMapping contain mapping between paginationtype and operator.
	PaginationTypeQueryMapping = map[string]Operator{
		"next":     LTOp,
		"previous": GTOp,
	}
	SortingTypeQueryMapping = map[string]SortDirection{
		"asc":  Asc,
		"desc": Desc,
	}
)

// ToMongo converts Filter object to mongo query and sets sorting and pagination options
// like sort, limit and skip.
func (f *QueryBuilder) ToMongo(opts *options.FindOptions) (bson.M, error) {

	var filter bson.M
	if query, ok := f.Query.(bson.M); !ok {
		return nil, fmt.Errorf("error reading query")
	} else {
		filter = query
	}

	// Make sort query mongo compatible
	if f.Sort != nil {
		sort := make(bson.D, 0, len(f.Sort))
		for _, v := range f.Sort {
			sort = append(sort, bson.E{Key: v.Name, Value: v.Direction})
		}
		opts.SetSort(sort)
	}

	if f.Limit > 0 {
		opts.SetLimit(f.Limit)

	}
	if f.Skip > 0 {
		opts.SetSkip(f.Skip)
	}

	return filter, nil
}

func GetAggregationQuery(filter bson.M, options *options.FindOptions) []bson.M {
	query := []bson.M{{"$match": filter}}

	if options.Sort != nil {
		query = append(query, bson.M{"$sort": options.Sort})
	}

	if options.Skip != nil {
		query = append(query, bson.M{"$skip": options.Skip})
	}

	if options.Limit != nil {
		query = append(query, bson.M{"$limit": options.Limit})
	}
	return query
}

// BuildQuery builds query using pagination params and filter
func BuildQuery(filter bson.M, paginationParameter paginator.PaginationQueryParam) (*QueryBuilder, error) {
	var qb QueryBuilder
	var skipCount int64
	// var paginationQ PaginationQ
	if !stringops.IsBlank(paginationParameter.LastID) {
		objectID, e := primitive.ObjectIDFromHex(paginationParameter.LastID)
		if e != nil {
			return nil, ErrorUnableToParseLastID
		}
		skipCount = 0
		operator := PaginationTypeQueryMapping[string(paginationParameter.Type)]
		filter["_id"] = bson.M{
			string(operator): objectID,
		}
	} else {
		skipCount = getSkipCount(paginationParameter)
	}
	qb.Query = filter
	qb.Skip = skipCount
	qb.Limit = paginationParameter.PageSize

	sortingFields, errFields := paginationParameter.GetSortingFields()
	if errFields != nil {
		return nil, errFields
	}

	for _, v := range sortingFields {
		qb.Sort = append(qb.Sort, SortType{
			Name:      v.Key,
			Direction: getSortTypeFromOrderBy(v.OrderBy),
		})
	}

	if len(qb.Sort) == 0 {
		qb.Sort = []SortType{{Name: "_id", Direction: Desc}}
	}

	return &qb, nil
}

// getSortTypeFromOrderBy returns sort type from order by
func getSortTypeFromOrderBy(orderBy paginator.OrderDirection) SortDirection {
	switch orderBy {
	case paginator.OrderByAsc:
		return Asc
	case paginator.OrderByDesc:
		return Desc
	}
	return None
}

// getSkipCount returns skip count from pagination params
func getSkipCount(paginationParameter paginator.PaginationQueryParam) int64 {
	if paginationParameter.PageNo == 0 {
		return 0
	}
	return (paginationParameter.PageNo - 1) * paginationParameter.PageSize
}
