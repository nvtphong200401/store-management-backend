package models

type JoinRequest struct {
	EmployeeID uint       `gorm:"primaryKey"`
	StoreID    uint       `gorm:"primaryKey"`
	Employee   Employee   `gorm:"foreignKey:EmployeeID"`
	Store      StoreModel `gorm:"foreignKey:StoreID" json:"-"`
	Status     StatusEnum `gorm:"default:'pending'"`
}
type StatusEnum string

const (
	PendingStatus  StatusEnum = "pending"
	AcceptedStatus StatusEnum = "accepted"
	DeniedStatus   StatusEnum = "denied"
)

func (jr *JoinRequest) IsPending() bool {
	return jr.Status == PendingStatus
}
