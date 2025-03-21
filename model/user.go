package model

import (
	"time"
)

type User struct {
	ID                   uint      `gorm:"primarykey"`
	FirstName            string    `gorm:"not null"`
	LastName             string    `gorm:"not null"`
	Email                string    `gorm:"uniqueIndex;not null"`
	GmailAccessToken     string    `gorm:"not null"`
	GmailRefreshToken    string    `gorm:"not null"`
	GmailTokenExpiry     time.Time `gorm:"not null"`
	IsGmailTokenValid    bool
	Picture              string
	LastScannedTimestamp time.Time `gorm:"index"`
	LastProcessedMail    string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}
