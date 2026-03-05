package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserMFA struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	TOTPSecret string    `json:"-" gorm:"not null"`
	IsEnabled  bool      `json:"is_enabled" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	User       User      `json:"-" gorm:"foreignKey:UserID"`
}

type RecoveryCode struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"index;not null"`
	CodeHash  string     `json:"-" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
	User      User       `json:"-" gorm:"foreignKey:UserID"`
}

type UserCompensation struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"index;not null"`
	BaseSalary    float64   `json:"base_salary" gorm:"not null"`
	EffectiveDate time.Time `json:"effective_date" gorm:"not null"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	User          User      `json:"-" gorm:"foreignKey:UserID"`
}

type MFAChallenge struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"-" gorm:"index;not null"`
	ChallengeToken string    `json:"-" gorm:"uniqueIndex;not null"`
	FailedAttempts int       `json:"-" gorm:"default:0"`
	ExpiresAt      time.Time `json:"-" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	User           User      `json:"-" gorm:"foreignKey:UserID"`
}
