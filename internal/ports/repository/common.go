package repository

// Common repository types and interfaces used across multiple repositories
// Note: Pagination has been moved to pagination.go

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
