package models

import (
	"strings"

	"github.com/benk-techworld/www-backend/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f Filters) SortField() string {

	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.Trim(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) SortDirection() int {
	if strings.HasPrefix(f.Sort, "-") {
		return -1
	}
	return 1
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilters(v *validator.Validator, filters Filters) {

	// Check that the page and page_size parameters contain sensible values.
	v.Check(filters.Page > 0, "page", "must be greater than zero")
	v.Check(filters.Page <= 100000, "page", "must be a maximum of 100 thousands")
	v.Check(filters.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(filters.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.

	v.Check(validator.PermittedValues(filters.Sort, filters.SortSafeList...), "sort", "invalid sort value")
}
