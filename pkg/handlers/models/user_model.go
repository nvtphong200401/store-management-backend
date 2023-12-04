package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                string `gorm:"unique"`
	Password             []byte `json:"-"`
	VerificationCode     string `json:"-"`
	Verified             bool   `json:"Verified" gorm:"default:false"`
	VerificationSentAt   time.Time
	VerificationAttempts int `json:"VerificationAttempts" gorm:"default:0"`
}
