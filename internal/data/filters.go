package data

import (
	"strings"

	"github.com/khatibomar/gogive/internal/validator"
)

type Filters struct {
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	PageSize     int     `json:"page_size,omitempty"`
	TotalRecords int     `json:"total_records,omitempty"`
	Cursor       *Cursor `json:"cursor,omitempty"`
}

type Cursor struct {
	LastSortVal interface{} `json:"last_sort_val,omitempty"`
	LastID      int64       `json:"last_id,omitempty"`
}

func newMetadata(totalRecords, pageSize int, cursor *Cursor) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		Cursor:       cursor,
	}
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func ValidateCursors(v *validator.Validator, c Cursor) {
	v.Check(c.LastID > 0, "last_id", "must be greater than zero")
	id, ok := c.LastSortVal.(int64)
	if ok {
		v.Check(id > 0, "last_sort_val", "must be greater than zero")
	}
}
