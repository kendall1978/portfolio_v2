package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

func ListRegions(c *gin.Context) {
	var regions []models.Region
	query := database.DB

	if areaType := c.Query("area_type"); areaType != "" {
		switch areaType {
		case "state":
			query = query.Where("area_type = ?", 2)
		case "metro":
			query = query.Where("area_type = ?", 4)
		}
	}

	query.Order("state, area_title").Find(&regions)
	c.JSON(http.StatusOK, regions)
}

func ListOccupations(c *gin.Context) {
	type Occupation struct {
		OccCode  string `json:"occ_code"`
		OccTitle string `json:"occ_title"`
	}

	var occupations []Occupation
	database.DB.Model(&models.SalarySnapshot{}).
		Select("DISTINCT occ_code, occ_title").
		Order("occ_title").
		Find(&occupations)

	c.JSON(http.StatusOK, occupations)
}

func ListSalaries(c *gin.Context) {
	var snapshots []models.SalarySnapshot
	query := database.DB.Preload("Region")

	if states := c.Query("states"); states != "" {
		codes := strings.Split(states, ",")
		query = query.Joins("JOIN regions ON regions.id = salary_snapshots.region_id").
			Where("regions.state_code IN ?", codes)
	}
	if occCodes := c.Query("occ_codes"); occCodes != "" {
		codes := strings.Split(occCodes, ",")
		query = query.Where("occ_code IN ?", codes)
	}

	query.Order("year desc").Limit(1000).Find(&snapshots)
	c.JSON(http.StatusOK, snapshots)
}

func SalaryTrends(c *gin.Context) {
	type TrendPoint struct {
		StateCode string  `json:"state_code"`
		State     string  `json:"state"`
		Year      int     `json:"year"`
		AvgMedian float64 `json:"avg_median"`
		AvgP25    float64 `json:"avg_p25"`
		AvgP75    float64 `json:"avg_p75"`
	}

	occCodes := c.Query("occ_codes")
	if occCodes == "" {
		c.JSON(http.StatusOK, []TrendPoint{})
		return
	}

	areaTypeVal := 2
	if areaType := c.Query("area_type"); areaType == "metro" {
		areaTypeVal = 4
	}

	query := `
		SELECT r.state_code, r.state,
			ss.year,
			AVG(ss.median_salary) as avg_median,
			AVG(ss.percentile25) as avg_p25,
			AVG(ss.percentile75) as avg_p75
		FROM salary_snapshots ss
		JOIN regions r ON r.id = ss.region_id
		WHERE r.area_type = ?
		AND ss.occ_code IN ?
	`
	args := []interface{}{areaTypeVal, strings.Split(occCodes, ",")}

	if states := c.Query("states"); states != "" {
		codes := strings.Split(states, ",")
		query += " AND r.state_code IN ?"
		args = append(args, codes)
	}

	query += `
		GROUP BY r.state_code, r.state, ss.year
		ORDER BY ss.year, r.state_code
	`

	var trends []TrendPoint
	database.DB.Raw(query, args...).Scan(&trends)
	c.JSON(http.StatusOK, trends)
}

func SalaryCompare(c *gin.Context) {
	userID := c.GetUint("user_id")

	var comp models.UserCompensation
	database.DB.Where("user_id = ?", userID).Order("effective_date desc").First(&comp)

	type StateComparison struct {
		StateCode string  `json:"state_code"`
		State     string  `json:"state"`
		Median    float64 `json:"median"`
		P25       float64 `json:"p25"`
		P75       float64 `json:"p75"`
	}

	occCodes := c.Query("occ_codes")
	if occCodes == "" {
		c.JSON(http.StatusOK, gin.H{
			"user_salary": comp.BaseSalary,
			"comparisons": []StateComparison{},
		})
		return
	}

	areaTypeVal := 2
	if areaType := c.Query("area_type"); areaType == "metro" {
		areaTypeVal = 4
	}

	query := `
		SELECT r.state_code, r.state,
			AVG(ss.median_salary) as median,
			AVG(ss.percentile25) as p25,
			AVG(ss.percentile75) as p75
		FROM salary_snapshots ss
		JOIN regions r ON r.id = ss.region_id
		WHERE r.area_type = ?
		AND ss.occ_code IN ?
		AND ss.year = (SELECT MAX(year) FROM salary_snapshots WHERE occ_code IN ?)
	`
	occList := strings.Split(occCodes, ",")
	args := []interface{}{areaTypeVal, occList, occList}

	if states := c.Query("states"); states != "" {
		codes := strings.Split(states, ",")
		query += " AND r.state_code IN ?"
		args = append(args, codes)
	}

	query += `
		GROUP BY r.state_code, r.state
		ORDER BY r.state_code
	`

	var comparisons []StateComparison
	database.DB.Raw(query, args...).Scan(&comparisons)

	c.JSON(http.StatusOK, gin.H{
		"user_salary": comp.BaseSalary,
		"comparisons": comparisons,
	})
}