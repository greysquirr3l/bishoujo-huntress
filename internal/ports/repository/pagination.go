package repository

// Pagination represents pagination information from API responses
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// Next returns the next page number, or the current page if already at the last page
func (p *Pagination) Next() int {
	if p.Page < p.TotalPages {
		return p.Page + 1
	}
	return p.Page
}

// Previous returns the previous page number, or the current page if already at the first page
func (p *Pagination) Previous() int {
	if p.Page > 1 {
		return p.Page - 1
	}
	return p.Page
}

// HasNext returns true if there is a next page
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

// HasPrevious returns true if there is a previous page
func (p *Pagination) HasPrevious() bool {
	return p.Page > 1
}

// IsFirstPage returns true if this is the first page
func (p *Pagination) IsFirstPage() bool {
	return p.Page == 1
}

// IsLastPage returns true if this is the last page
func (p *Pagination) IsLastPage() bool {
	return p.Page == p.TotalPages
}
