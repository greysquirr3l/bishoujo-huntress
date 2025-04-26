package huntress_test

import (
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestPaginationStruct(t *testing.T) {
	p := huntress.Pagination{
		CurrentPage: 1,
		PerPage:     10,
		TotalPages:  2,
		TotalItems:  20,
	}
	if p.CurrentPage != 1 || p.PerPage != 10 || p.TotalPages != 2 || p.TotalItems != 20 {
		t.Error("pagination struct fields not set correctly")
	}
}
