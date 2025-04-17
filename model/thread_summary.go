package model

import (
	"time"

	"github.com/lib/pq"
)

type ThreadSummary struct {
	ID                       uint           `gorm:"primarykey" json:"id"`
	GmailThreadId            string         `gorm:"uniqueIndex;not null" json:"gmailThreadID"`
	ProcessedEmailIds        pq.StringArray `gorm:"type:text[]" json:"-"`
	Recipients               pq.StringArray `gorm:"type:text[]" json:"recipients"`
	Summary                  string         `json:"summary"`
	ThreadSubject            string         `json:"threadSubject"`
	UserID                   uint           `gorm:"index;not null" json:"-"`
	User                     User           `json:"-"`
	MostRecentEmailTimestamp time.Time      `json:"mostRecentEmailTimestamp"`

	CreatedAt time.Time `gorm:"index;not null" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"-"`
}
