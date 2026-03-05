package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

type CompensationRequest struct {
	BaseSalary    float64 `json:"base_salary" binding:"required"`
	EffectiveDate string  `json:"effective_date" binding:"required"`
	Notes         string  `json:"notes"`
}

func ListCompensation(c *gin.Context) {
	userID := c.GetUint("user_id")
	var entries []models.UserCompensation
	database.DB.Where("user_id = ?", userID).Order("effective_date desc").Find(&entries)
	c.JSON(http.StatusOK, entries)
}

func CreateCompensation(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CompensationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry := models.UserCompensation{
		UserID:     userID,
		BaseSalary: req.BaseSalary,
		Notes:      req.Notes,
	}

	if t, err := time.Parse(time.RFC3339, req.EffectiveDate); err == nil {
		entry.EffectiveDate = t
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use RFC3339"})
		return
	}

	if result := database.DB.Create(&entry); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func UpdateCompensation(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compensation ID"})
		return
	}

	var entry models.UserCompensation
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	var req CompensationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.BaseSalary = req.BaseSalary
	entry.Notes = req.Notes
	t, err := time.Parse(time.RFC3339, req.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use RFC3339"})
		return
	}
	entry.EffectiveDate = t

	database.DB.Save(&entry)
	c.JSON(http.StatusOK, entry)
}

func DeleteCompensation(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compensation ID"})
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserCompensation{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
