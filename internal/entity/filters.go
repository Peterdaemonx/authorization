package entity

import "strings"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	return ""
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	// ASC is the default value of the DB. By returning an empty value instead of ASC
	// we can return unsorted result sets as well. Order by can be a performance bottleneck
	return ""
}
