package types

// Pagination represents pagination information for API responses
type Pagination struct {
	// Page represents the current page number
	Page int `json:"page"`
	// PerPage represents the number of items per page
	PerPage int `json:"per_page"`
	// TotalItems represents the total number of items across all pages
	TotalItems int `json:"total_items"`
	// TotalPages represents the total number of pages
	TotalPages int `json:"total_pages"`
}

// HasNext returns true if there are more pages available
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

// HasPrevious returns true if there are previous pages available
func (p *Pagination) HasPrevious() bool {
	return p.Page > 1
}

// NextPage returns the next page number, or the current page if already on the last page
func (p *Pagination) NextPage() int {
	if p.HasNext() {
		return p.Page + 1
	}
	return p.Page
}

// PreviousPage returns the previous page number, or the current page if already on the first page
func (p *Pagination) PreviousPage() int {
	if p.HasPrevious() {
		return p.Page - 1
	}
	return p.Page
}

// IsFirstPage returns true if this is the first page
func (p *Pagination) IsFirstPage() bool {
	return p.Page == 1
}

// IsLastPage returns true if this is the last page
func (p *Pagination) IsLastPage() bool {
	return p.Page == p.TotalPages
}
