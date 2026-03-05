package models

import "time"

type SalarySnapshot struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	RegionID        uint      `json:"region_id" gorm:"not null;uniqueIndex:idx_region_occ_year"`
	Source          string    `json:"source" gorm:"not null"`
	OccCode         string    `json:"occ_code" gorm:"not null;uniqueIndex:idx_region_occ_year;index"`
	OccTitle        string    `json:"occ_title" gorm:"not null"`
	TotalEmployment int       `json:"total_employment"`
	MedianSalary    float64   `json:"median_salary" gorm:"not null"`
	Percentile25    float64   `json:"percentile_25"`
	Percentile75    float64   `json:"percentile_75"`
	Percentile10    float64   `json:"percentile_10"`
	Percentile90    float64   `json:"percentile_90"`
	MeanSalary      float64   `json:"mean_salary"`
	Year            int       `json:"year" gorm:"not null;uniqueIndex:idx_region_occ_year;index"`
	SnapshotDate    time.Time `json:"snapshot_date" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at"`
	Region          Region    `json:"region" gorm:"foreignKey:RegionID"`
}