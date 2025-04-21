package repository

// Pagination contains information about pagination
type Pagination struct {
	Page       int
	Limit      int
	TotalItems int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

// NewPagination creates a new pagination object
func NewPagination(page, limit, totalItems int) Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := (totalItems + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// OrderDirection defines the direction of ordering
type OrderDirection string

const (
	// OrderAsc indicates ascending order
	OrderAsc OrderDirection = "asc"
	// OrderDesc indicates descending order
	OrderDesc OrderDirection = "desc"
)

// OrderBy defines an ordering specification
type OrderBy struct {
	Field     string
	Direction OrderDirection
}

// TimeRange defines a time range for filtering
type TimeRange struct {
	Start string
	End   string
}
