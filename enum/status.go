package enum

// Status represents the status of a resource
type Status string

const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
	Present  Status = "present"
	Absent   Status = "absent"
	Active   Status = "active"
	Inactive Status = "inactive"
)

// Statuses returns all the available statuses
func Statuses() []Status {
	return []Status{
		Pending,
		Accepted,
		Rejected,
		Present,
		Absent,
		Active,
		Inactive,
	}
}

// IsValidStatus checks if the given status is valid
func IsValidStatus(status Status) bool {
	for _, s := range Statuses() {
		if s == status {
			return true
		}
	}
	return false
}

// IsPending checks if the given status is pending
func IsPending(status Status) bool {
	return status == Pending
}

// IsAccepted checks if the given status is accepted
func IsAccepted(status Status) bool {
	return status == Accepted
}

// IsRejected checks if the given status is rejected
func IsRejected(status Status) bool {
	return status == Rejected
}

// IsPresent checks if the given status is present
func IsPresent(status Status) bool {
	return status == Present
}

// IsAbsent checks if the given status is absent
func IsAbsent(status Status) bool {
	return status == Absent
}

// IsActive checks if the given status is active
func IsActive(status Status) bool {
	return status == Active
}

// IsInactive checks if the given status is inactive
func IsInactive(status Status) bool {
	return status == Inactive
}
