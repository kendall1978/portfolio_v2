package models

type Region struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	AreaCode  string `json:"area_code" gorm:"uniqueIndex;not null"`
	AreaTitle string `json:"area_title" gorm:"not null"`
	AreaType  int    `json:"area_type" gorm:"not null;index"`
	State     string `json:"state" gorm:"not null;index"`
	StateCode string `json:"state_code" gorm:"not null;size:2;index"`
}