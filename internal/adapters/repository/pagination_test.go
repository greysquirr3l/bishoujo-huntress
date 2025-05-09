package repository

import (
	"testing"

	repo "github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

func TestPagination_Next(t *testing.T) {
	p := repo.Pagination{Page: 1, TotalPages: 3}
	if next := p.Next(); next != 2 {
		t.Errorf("expected next=2, got %d", next)
	}
	p.Page = 3
	if next := p.Next(); next != 3 {
		t.Errorf("expected next=3 (last page), got %d", next)
	}
}

func TestPagination_Previous(t *testing.T) {
	p := repo.Pagination{Page: 2, TotalPages: 3}
	if prev := p.Previous(); prev != 1 {
		t.Errorf("expected prev=1, got %d", prev)
	}
	p.Page = 1
	if prev := p.Previous(); prev != 1 {
		t.Errorf("expected prev=1 (first page), got %d", prev)
	}
}

func TestPagination_HasNextAndPrevious(t *testing.T) {
	p := repo.Pagination{Page: 1, TotalPages: 2}
	if !p.HasNext() {
		t.Errorf("expected HasNext=true")
	}
	if p.HasPrevious() {
		t.Errorf("expected HasPrevious=false")
	}
	p.Page = 2
	if p.HasNext() {
		t.Errorf("expected HasNext=false on last page")
	}
	if !p.HasPrevious() {
		t.Errorf("expected HasPrevious=true on page 2")
	}
}

func TestPagination_IsFirstAndLastPage(t *testing.T) {
	p := repo.Pagination{Page: 1, TotalPages: 2}
	if !p.IsFirstPage() {
		t.Errorf("expected IsFirstPage=true")
	}
	if p.IsLastPage() {
		t.Errorf("expected IsLastPage=false")
	}
	p.Page = 2
	if p.IsFirstPage() {
		t.Errorf("expected IsFirstPage=false on page 2")
	}
	if !p.IsLastPage() {
		t.Errorf("expected IsLastPage=true on last page")
	}
}
