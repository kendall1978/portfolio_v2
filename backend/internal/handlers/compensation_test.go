package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/handlers"
	"github.com/kendall1978/project_salary/backend/internal/middleware"
)

func setupCompRouter() *gin.Engine {
	router := setupRouter()
	comp := router.Group("/api/compensation")
	comp.Use(middleware.AuthRequired())
	{
		comp.GET("", handlers.ListCompensation)
		comp.POST("", handlers.CreateCompensation)
		comp.PUT("/:id", handlers.UpdateCompensation)
		comp.DELETE("/:id", handlers.DeleteCompensation)
	}
	return router
}

func TestCreateCompensation_Success(t *testing.T) {
	setupTestDB(t)
	router := setupCompRouter()
	token := registerAndLogin(t, router)

	body := map[string]interface{}{
		"base_salary":    85000.00,
		"effective_date": "2026-01-01T00:00:00Z",
		"notes":          "Current salary",
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/compensation", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["base_salary"] != 85000.0 {
		t.Errorf("Expected base_salary 85000, got %v", response["base_salary"])
	}
}

func TestListCompensation_Success(t *testing.T) {
	setupTestDB(t)
	router := setupCompRouter()
	token := registerAndLogin(t, router)

	// Create one first
	body := map[string]interface{}{
		"base_salary":    85000.00,
		"effective_date": "2026-01-01T00:00:00Z",
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/compensation", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// List
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/compensation", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(response))
	}
}

func TestUpdateCompensation_Success(t *testing.T) {
	setupTestDB(t)
	router := setupCompRouter()
	token := registerAndLogin(t, router)

	// Create
	body := map[string]interface{}{
		"base_salary":    85000.00,
		"effective_date": "2026-01-01T00:00:00Z",
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/compensation", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := fmt.Sprintf("%.0f", created["id"].(float64))

	// Update
	body["base_salary"] = 95000.00
	jsonBody, _ = json.Marshal(body)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/compensation/"+id, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated["base_salary"] != 95000.0 {
		t.Errorf("Expected base_salary 95000, got %v", updated["base_salary"])
	}
}

func TestDeleteCompensation_Success(t *testing.T) {
	setupTestDB(t)
	router := setupCompRouter()
	token := registerAndLogin(t, router)

	// Create
	body := map[string]interface{}{
		"base_salary":    85000.00,
		"effective_date": "2026-01-01T00:00:00Z",
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/compensation", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := fmt.Sprintf("%.0f", created["id"].(float64))

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/compensation/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify it's gone
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/compensation", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var list []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Errorf("Expected 0 entries after delete, got %d", len(list))
	}
}
