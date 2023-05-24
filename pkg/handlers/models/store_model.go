package models

import (
	"gorm.io/gorm"
)

type StoreModel struct {
	gorm.Model
	StoreName string `json:"StoreName"`
	Address   string `json:"Address"`
}
