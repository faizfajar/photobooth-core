package domain

type UserRole string

const (
	RoleOwner UserRole = "owner"
	RoleStaff UserRole = "staff"
)

// status booth
type BoothStatus string

const (
	BoothActive      BoothStatus = "active"
	BoothMaintenance BoothStatus = "maintenance"
	BoothOffline     BoothStatus = "offline"
)

// status transaction
type TransactionStatus string

const (
	TransPending   TransactionStatus = "pending"
	TransCompleted TransactionStatus = "completed"
	TransFailed    TransactionStatus = "failed"
)
