package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/handlers"
	"github.com/kendall1978/project_salary/backend/internal/middleware"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

func TestListRegions(t *testing.T) {
	setupTestDB(t)

	// Create a test region with the new schema
	database.DB.Create(&models.Region{
		AreaCode:  "01",
		AreaTitle: "Alabama",
		AreaType:  2,
		State:     "Alabama",
		StateCode: "AL",
	})

	router := setupRouter()
	router.GET("/api/regions", middleware.AuthRequired(), handlers.ListRegions)

	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/regions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var regions []models.Region
	json.Unmarshal(w.Body.Bytes(), &regions)
	if len(regions) == 0 {
		t.Error("Expected regions, got none")
	}
}

func TestListSalaries_Empty(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	router.GET("/api/salaries", middleware.AuthRequired(), handlers.ListSalaries)

	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/salaries", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var snapshots []models.SalarySnapshot
	json.Unmarshal(w.Body.Bytes(), &snapshots)
	if len(snapshots) != 0 {
		t.Errorf("Expected 0 snapshots, got %d", len(snapshots))
	}
}

func TestSalaryCompare_NoData(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	router.GET("/api/salaries/compare", middleware.AuthRequired(), handlers.SalaryCompare)

	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/salaries/compare", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["user_salary"] != 0.0 {
		t.Errorf("Expected user_salary 0, got %v", response["user_salary"])
	}
}