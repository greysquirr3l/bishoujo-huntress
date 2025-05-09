package common

import "testing"

func TestPagination_HasNextAndPrevious(t *testing.T) {
	p := &Pagination{Page: 2, TotalPages: 3}
	if !p.HasNext() {
		t.Error("expected HasNext true")
	}
	if !p.HasPrevious() {
		t.Error("expected HasPrevious true")
	}
}

func TestPagination_NextPage_PreviousPage(t *testing.T) {
	p := &Pagination{Page: 1, TotalPages: 2}
	if p.NextPage() != 2 {
		t.Error("expected next page 2")
	}
	if p.PreviousPage() != 1 {
		t.Error("expected previous page 1")
	}
	p.Page = 2
	if p.NextPage() != 2 {
		t.Error("expected next page 2 (last)")
	}
	if p.PreviousPage() != 1 {
		t.Error("expected previous page 1")
	}
}

func TestPagination_IsFirstAndLastPage(t *testing.T) {
	p := &Pagination{Page: 1, TotalPages: 1}
	if !p.IsFirstPage() || !p.IsLastPage() {
		t.Error("expected first and last page true")
	}
	p.TotalPages = 2
	p.Page = 2
	if p.IsFirstPage() {
		t.Error("expected not first page")
	}
	if !p.IsLastPage() {
		t.Error("expected last page true")
	}
}

func TestNewAndDefaultPagination(t *testing.T) {
	p := NewPagination(2, 10, 100, 10)
	if p.Page != 2 || p.PerPage != 10 || p.TotalItems != 100 || p.TotalPages != 10 {
		t.Error("unexpected pagination fields")
	}
	def := DefaultPagination()
	if def.Page != 1 || def.PerPage != 100 {
		t.Error("unexpected default pagination")
	}
}
