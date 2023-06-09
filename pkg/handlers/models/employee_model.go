package models

type Employee struct {
	User
	StoreID  uint             `json:"StoreID,omitempty"`
	Position EmployeePosition `gorm:"default:'unknown'"`
}

type EmployeePosition string

const (
	Owner   EmployeePosition = "owner"
	Staff   EmployeePosition = "staff"
	Unknown EmployeePosition = "unknown"
)

func (e *Employee) IsEmployeeOwner() bool {
	return e.Position == Owner
}
func (employee *Employee) AlreadyInStore() bool {
	return employee.StoreID != 0
}
