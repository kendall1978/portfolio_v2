package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

type BlogPostRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	ImagePath string `json:"image_path"`
	Date      string `json:"date" binding:"required"`
}

// ListBlogPosts - PUBLIC, no auth required
func ListBlogPosts(c *gin.Context) {
	var posts []models.BlogPost
	database.DB.Order("date desc").Find(&posts)
	c.JSON(http.StatusOK, posts)
}

// GetBlogPost - PUBLIC, no auth required
func GetBlogPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.BlogPost
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// CreateBlogPost - AUTH required
func CreateBlogPost(c *gin.Context) {
	var req BlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.BlogPost{
		Title:     req.Title,
		Content:   req.Content,
		ImagePath: req.ImagePath,
	}

	if t, err := time.Parse("2006-01-02", req.Date); err == nil {
		post.Date = t
	} else if t, err := time.Parse(time.RFC3339, req.Date); err == nil {
		post.Date = t
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD or RFC3339"})
		return
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdateBlogPost - AUTH required
func UpdateBlogPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.BlogPost
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var req BlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.ImagePath = req.ImagePath

	if t, err := time.Parse("2006-01-02", req.Date); err == nil {
		post.Date = t
	} else if t, err := time.Parse(time.RFC3339, req.Date); err == nil {
		post.Date = t
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	database.DB.Save(&post)
	c.JSON(http.StatusOK, post)
}

// DeleteBlogPost - AUTH required
func DeleteBlogPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	result := database.DB.Delete(&models.BlogPost{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
