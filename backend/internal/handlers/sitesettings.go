package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

// GetSiteSettings - PUBLIC
func GetSiteSettings(c *gin.Context) {
	var settings []models.SiteSetting
	database.DB.Find(&settings)

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, result)
}

// UpdateSiteSettings - AUTH required
func UpdateSiteSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range req {
		var setting models.SiteSetting
		result := database.DB.Where("key = ?", key).First(&setting)
		if result.Error != nil {
			setting = models.SiteSetting{Key: key, Value: value}
			database.DB.Create(&setting)
		} else {
			setting.Value = value
			database.DB.Save(&setting)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}
