package domain

type Decision string

const (
	DecisionSafeToDelete Decision = "safe_to_delete"
	DecisionPending      Decision = "pending"
)
