package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kendall1978/project_salary/backend/internal/handlers"
	"github.com/kendall1978/project_salary/backend/internal/middleware"
)

func TestGetProfile_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	router.GET("/api/profile", middleware.AuthRequired(), handlers.GetProfile)

	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["first_name"] != "Kendall" {
		t.Errorf("Expected first_name Kendall, got %v", response["first_name"])
	}
	if response["last_name"] != "Roberts" {
		t.Errorf("Expected last_name Roberts, got %v", response["last_name"])
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	router.PUT("/api/profile", middleware.AuthRequired(), handlers.UpdateProfile)
	router.GET("/api/profile", middleware.AuthRequired(), handlers.GetProfile)

	token := registerAndLogin(t, router)

	// Update
	body := map[string]string{"first_name": "Kenny", "last_name": "Roberts"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify the update persisted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["first_name"] != "Kenny" {
		t.Errorf("Expected first_name Kenny, got %v", response["first_name"])
	}
}

func TestGetProfile_Unauthorized(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	router.GET("/api/profile", middleware.AuthRequired(), handlers.GetProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profile", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}
