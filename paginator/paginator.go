package paginator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/a01k-io/modules/stringops"
	"github.com/gorilla/schema"
)

var (
	//URLParamDecoder returns a new Decoder
	//For parsing url query params
	URLParamDecoder       = schema.NewDecoder()
	minPageSize     int64 = 50
	maxPageSize     int64 = 100
)

// QueryParamFilter defines settings for query parser
type QueryParamFilter struct {
	Name              string
	MaxLen            string
	IgnoreUnknownKeys bool
	Strip             bool
}

// Filter defines settings for parser
type Filter struct {
	DisallowUnknownFields bool
	MaxSize               int64
}

// NewQueryParamsFromReq parse the incoming http request query params payload into the given data type.
// i must be pointer to the data type.
func NewQueryParamsFromReq(i interface{}, r *http.Request, filter QueryParamFilter) error {
	if err := r.ParseForm(); err != nil {
		return errors.New("failed to parse query param")
	}
	if filter.IgnoreUnknownKeys {
		URLParamDecoder.IgnoreUnknownKeys(true)
	}

	if err := URLParamDecoder.Decode(i, r.Form); err != nil {
		return errors.New("failed to Decode query param")
	}
	return nil
}

// MultipartFormParamFilter defines settings for form parser
type MultipartFormParamFilter struct {
	MaxMemory         int64
	IgnoreUnknownKeys bool
}

// NewMultipartFormParamsFromReq parse the incoming http request multipart form params payload into the given data type.
// will only used to parse non-binary form data. for parsing Binary data like files should use r.FormFile
// i must be pointer to the data type.
func NewMultipartFormParamsFromReq(i interface{}, r *http.Request, filter MultipartFormParamFilter) *fiber.Error {
	var defaultMultipartFormMaxMemory int64 = 32 << 20 // 32 MB
	if filter.MaxMemory != 0 {
		defaultMultipartFormMaxMemory = filter.MaxMemory
	}
	if e := r.ParseMultipartForm(defaultMultipartFormMaxMemory); e != nil {
		return fiber.NewError(http.StatusBadRequest, "failed to parse multipart form data.")
	}
	if filter.IgnoreUnknownKeys {
		URLParamDecoder.IgnoreUnknownKeys(true)
	}

	if e := URLParamDecoder.Decode(i, r.PostForm); e != nil {
		return fiber.NewError(http.StatusBadRequest, "failed to Decode multipart data")
	}
	return nil
}

// PaginationQueryType defines type of paginated query
type PaginationQueryType string

const (
	//NextPage will return next records
	NextPage PaginationQueryType = "next"
	//PrevPage will return previous record
	PrevPage PaginationQueryType = "prev"
)

// Valid validate paginated query type
func (t PaginationQueryType) Valid() bool {
	switch t {
	case PrevPage, NextPage:
		return true
	}
	return false
}

// OrderDirection defines type of order_by query param
type OrderDirection string

const (
	//OrderByAsc will sort document by ascending order
	OrderByAsc OrderDirection = "asc"
	//OrderByDesc will sort document by descending order
	OrderByDesc OrderDirection = "desc"
)

// Valid validate paginated order_by query type
func (t OrderDirection) Valid() bool {
	return t == OrderByAsc || t == OrderByDesc
}

// SortingField define schema for sorting query params
type SortingField struct {
	Key     string
	OrderBy OrderDirection
}

// GetSortingFields convert sorting query params to SortingField and validate it
func getSortingFields(params []string) ([]*SortingField, error) {
	var errFields []string
	sortingFields := make([]*SortingField, 0, len(params))

	for _, s := range params {
		sortField := strings.Split(s, ":")
		if len(sortField) != 2 {
			errFields = append(errFields, fmt.Sprintf("invalid sort_by value: %v", s))
			continue
		}

		if stringops.IsBlank(sortField[0]) {
			errFields = append(errFields, "sort_by field name is required")
		}

		if !OrderDirection(sortField[1]).Valid() {
			errFields = append(errFields, fmt.Sprintf("invalid sort_by order value: %v", s))
		}

		sortingFields = append(sortingFields, &SortingField{Key: sortField[0], OrderBy: OrderDirection(sortField[1])})
	}

	if len(errFields) != 0 {
		return make([]*SortingField, 0), errors.Wrapf(errors.New("invalid sort_by query params: "), fmt.Sprintf("%v", errFields))
	}

	return sortingFields, nil
}

// PaginationQueryParam define schema for pagination query params
type PaginationQueryParam struct {
	PageNo   int64               `schema:"page_no" query:"page_no" json:"page_no"`
	PageSize int64               `schema:"page_size" query:"page_size" json:"page_size"`
	LastID   string              `schema:"last_id" query:"last_id" json:"last_id"`
	Type     PaginationQueryType `schema:"pagination_type" query:"pagination_type" json:"pagination_type"`
	SortBy   []string            `schema:"sort_by" query:"sort_by" json:"sort_by"`
}

// GetSortingFields parse sort query param and return SortingField array
func (p *PaginationQueryParam) GetSortingFields() ([]*SortingField, error) {
	return getSortingFields(p.SortBy)
}

// Validate pagination query params
func (p *PaginationQueryParam) Validate() error {
	var errFields []string

	if !stringops.IsBlank(p.LastID) && !p.Type.Valid() {
		errFields = append(errFields, "Invalid pagination_type")
	} else if stringops.IsBlank(p.LastID) && p.PageNo <= 0 {
		errFields = append(errFields, fmt.Sprintf("invalid page_no value: %v", p.PageNo))
	}

	if p.PageSize < minPageSize || p.PageSize > maxPageSize {
		errFields = append(errFields, fmt.Sprintf("invalid page_size value it should be in between [%v - %v]", minPageSize, maxPageSize))
	}

	_, sortParamErr := p.GetSortingFields()
	if sortParamErr != nil {
		errFields = append(errFields, sortParamErr.Error())
	}
	if len(errFields) != 0 {
		return errors.Wrapf(errors.New("invalid pagination query params: "), fmt.Sprintf("%v", errFields))
	}

	return nil
}

// NewPaginationQueryParams parse the incoming http request to get pagination query params.
func NewPaginationQueryParams(r *http.Request) (*PaginationQueryParam, error) {
	var params PaginationQueryParam
	if err := NewQueryParamsFromReq(&params, r, QueryParamFilter{}); err != nil {
		return nil, err
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return &params, nil
}

// NewPaginationQueryParamsF parse the incoming http request to get pagination query params with filter
func NewPaginationQueryParamsF(r *http.Request, filter QueryParamFilter) (*PaginationQueryParam, error) {
	var params PaginationQueryParam
	if err := NewQueryParamsFromReq(&params, r, filter); err != nil {
		return nil, err
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return &params, nil
}

// PaginationInfo contain info about pagination
type PaginationInfo struct {
	PageNo     int64 `json:"page_no"`
	PageSize   int64 `json:"page_size"`
	TotalCount int   `json:"total_count"`
}

// PaginatedResponse contain records and pagination info
type PaginatedResponse struct {
	Records    interface{}    `json:"records"`
	Pagination PaginationInfo `json:"pagination"`
}

// CreatePaginatedAPIResponse add pagination info in api response
func CreatePaginatedAPIResponse(records interface{}, paginatedQueryParam PaginationQueryParam, totalCount int) PaginatedResponse {
	return PaginatedResponse{
		Records: records,
		Pagination: PaginationInfo{
			PageSize:   paginatedQueryParam.PageSize,
			PageNo:     paginatedQueryParam.PageNo,
			TotalCount: totalCount,
		},
	}
}

func trimArray(arrayList []string) {
	for k := range arrayList {
		arrayList[k] = strings.TrimSpace(arrayList[k])
	}
}
