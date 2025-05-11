package types

import "testing"

func TestPagination_HasNext(t *testing.T) {
	p := &Pagination{Page: 1, TotalPages: 3}
	if !p.HasNext() {
		t.Error("expected HasNext to be true when more pages exist")
	}
	p.Page = 3
	if p.HasNext() {
		t.Error("expected HasNext to be false on last page")
	}
}

func TestPagination_HasPrevious(t *testing.T) {
	p := &Pagination{Page: 2}
	if !p.HasPrevious() {
		t.Error("expected HasPrevious to be true when not on first page")
	}
	p.Page = 1
	if p.HasPrevious() {
		t.Error("expected HasPrevious to be false on first page")
	}
}

func TestPagination_NextPage(t *testing.T) {
	p := &Pagination{Page: 1, TotalPages: 2}
	if p.NextPage() != 2 {
		t.Errorf("expected NextPage to be 2, got %d", p.NextPage())
	}
	p.Page = 2
	if p.NextPage() != 2 {
		t.Errorf("expected NextPage to stay at last page, got %d", p.NextPage())
	}
}

func TestPagination_PreviousPage(t *testing.T) {
	p := &Pagination{Page: 2}
	if p.PreviousPage() != 1 {
		t.Errorf("expected PreviousPage to be 1, got %d", p.PreviousPage())
	}
	p.Page = 1
	if p.PreviousPage() != 1 {
		t.Errorf("expected PreviousPage to stay at first page, got %d", p.PreviousPage())
	}
}

func TestPagination_IsFirstPage(t *testing.T) {
	p := &Pagination{Page: 1}
	if !p.IsFirstPage() {
		t.Error("expected IsFirstPage to be true on first page")
	}
	p.Page = 2
	if p.IsFirstPage() {
		t.Error("expected IsFirstPage to be false on non-first page")
	}
}

func TestPagination_IsLastPage(t *testing.T) {
	p := &Pagination{Page: 3, TotalPages: 3}
	if !p.IsLastPage() {
		t.Error("expected IsLastPage to be true on last page")
	}
	p.Page = 2
	if p.IsLastPage() {
		t.Error("expected IsLastPage to be false on non-last page")
	}
}
