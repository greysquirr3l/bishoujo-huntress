package account

// Status represents the current status of an account
type Status string

const (
	// StatusActive indicates an active account
	StatusActive Status = "active"
	// StatusSuspended indicates a suspended account
	StatusSuspended Status = "suspended"
	// StatusDeactivated indicates a deactivated account
	StatusDeactivated Status = "deactivated"
	// StatusPending indicates a pending account
	StatusPending Status = "pending"
)

// IsValid checks if the status is a valid account status
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusSuspended, StatusDeactivated, StatusPending:
		return true
	default:
		return false
	}
}

// String returns the string representation of the account status
func (s Status) String() string {
	return string(s)
}
