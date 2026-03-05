package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
	"github.com/kendall1978/project_salary/backend/internal/scraper"
)

type ScrapeRequest struct {
	Backfill    bool   `json:"backfill"`
	LocalZipDir string `json:"local_zip_dir"`
}

func ScrapeTrigger(c *gin.Context) {
	var req ScrapeRequest
	c.ShouldBindJSON(&req)

	go func() {
		scraper.FetchOESData(req.Backfill, req.LocalZipDir)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "OES data import started",
	})
}

func ScrapeLogs(c *gin.Context) {
	var logs []models.ScrapeLog
	database.DB.Order("started_at desc").Limit(50).Find(&logs)
	c.JSON(http.StatusOK, logs)
}