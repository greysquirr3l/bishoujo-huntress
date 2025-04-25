package account

// AccountStatus represents the current status of an account
type AccountStatus string

const (
	// StatusActive indicates an active account
	StatusActive AccountStatus = "active"
	// StatusSuspended indicates a suspended account
	StatusSuspended AccountStatus = "suspended"
	// StatusDeactivated indicates a deactivated account
	StatusDeactivated AccountStatus = "deactivated"
	// StatusPending indicates a pending account
	StatusPending AccountStatus = "pending"
)

// IsValid checks if the status is a valid account status
func (s AccountStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusSuspended, StatusDeactivated, StatusPending:
		return true
	default:
		return false
	}
}

// String returns the string representation of the account status
func (s AccountStatus) String() string {
	return string(s)
}
