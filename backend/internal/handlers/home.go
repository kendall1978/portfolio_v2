package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

type HomeSectionRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// ListHomeSections - PUBLIC
func ListHomeSections(c *gin.Context) {
	var sections []models.HomeSection
	database.DB.Order("sort_order asc").Find(&sections)
	c.JSON(http.StatusOK, sections)
}

// CreateHomeSection - AUTH required
func CreateHomeSection(c *gin.Context) {
	var req HomeSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section := models.HomeSection{
		Title:     req.Title,
		Content:   req.Content,
		SortOrder: req.SortOrder,
	}

	if err := database.DB.Create(&section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create section"})
		return
	}
	c.JSON(http.StatusCreated, section)
}

// UpdateHomeSection - AUTH required
func UpdateHomeSection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var section models.HomeSection
	if err := database.DB.First(&section, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	var req HomeSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section.Title = req.Title
	section.Content = req.Content
	section.SortOrder = req.SortOrder

	database.DB.Save(&section)
	c.JSON(http.StatusOK, section)
}

// DeleteHomeSection - AUTH required
func DeleteHomeSection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	result := database.DB.Delete(&models.HomeSection{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
