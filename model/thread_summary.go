package model

import (
	"time"

	"github.com/lib/pq"
)

type ThreadSummary struct {
	ID                       uint           `gorm:"primarykey"`
	GmailThreadId            string         `gorm:"uniqueIndex;not null"`
	ProcessedEmailIds        pq.StringArray `gorm:"type:text[]"`
	Summary                  string
	ThreadSubject            string
	UserID                   uint `gorm:"index;not null"`
	User                     User
	MostRecentEmailTimestamp time.Time

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
	DeletedAt time.Time
}
