package models

import "time"

type ScrapeLog struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Source       string     `json:"source" gorm:"not null"`
	Status       string     `json:"status" gorm:"not null"`
	RecordsFound int        `json:"records_found"`
	ErrorMessage string     `json:"error_message"`
	StartedAt    time.Time  `json:"started_at" gorm:"not null"`
	CompletedAt  *time.Time `json:"completed_at"`
}
