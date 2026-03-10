package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

type SocialLinkRequest struct {
	Platform  string `json:"platform" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Icon      string `json:"icon" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// ListSocialLinks - PUBLIC
func ListSocialLinks(c *gin.Context) {
	var links []models.SocialLink
	database.DB.Order("sort_order asc").Find(&links)
	c.JSON(http.StatusOK, links)
}

// CreateSocialLink - AUTH required
func CreateSocialLink(c *gin.Context) {
	var req SocialLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link := models.SocialLink{
		Platform:  req.Platform,
		URL:       req.URL,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}

	if err := database.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		return
	}
	c.JSON(http.StatusCreated, link)
}

// UpdateSocialLink - AUTH required
func UpdateSocialLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}

	var link models.SocialLink
	if err := database.DB.First(&link, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	var req SocialLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link.Platform = req.Platform
	link.URL = req.URL
	link.Icon = req.Icon
	link.SortOrder = req.SortOrder

	database.DB.Save(&link)
	c.JSON(http.StatusOK, link)
}

// DeleteSocialLink - AUTH required
func DeleteSocialLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}

	result := database.DB.Delete(&models.SocialLink{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
