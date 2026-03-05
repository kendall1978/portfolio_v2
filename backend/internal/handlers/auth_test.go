package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/crypto"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/handlers"
	"github.com/kendall1978/project_salary/backend/internal/middleware"
	"github.com/kendall1978/project_salary/backend/internal/models"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	handlers.SetJWTSecret("test-secret")
	err := database.Connect("postgres://postgres:postgres@localhost:5432/kroberts_personal_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	database.DB.AutoMigrate(
		&models.User{},
		&models.UserMFA{},
		&models.RecoveryCode{},
		&models.UserCompensation{},
		&models.MFAChallenge{},
		&models.Region{},
		&models.SalarySnapshot{},
		&models.ScrapeLog{},
	)
	// Clean tables before each test
	database.DB.Exec("DELETE FROM salary_snapshots")
	database.DB.Exec("DELETE FROM scrape_logs")
	database.DB.Exec("DELETE FROM regions")
	database.DB.Exec("DELETE FROM mfa_challenges")
	database.DB.Exec("DELETE FROM user_compensations")
	database.DB.Exec("DELETE FROM recovery_codes")
	database.DB.Exec("DELETE FROM user_mfas")
	database.DB.Exec("DELETE FROM users")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/mfa/verify", handlers.MFAVerify)
	auth.POST("/mfa/recover", handlers.MFARecover)
	auth.POST("/mfa/setup", middleware.AuthRequired(), handlers.MFASetup)
	auth.POST("/refresh", middleware.AuthRequired(), handlers.RefreshToken)
	return r
}

// Helper: register a user and return the JWT token
func registerAndLogin(t *testing.T, router *gin.Engine) string {
	t.Helper()
	regBody := map[string]string{
		"email": "test@example.com", "password": "securepassword123",
		"first_name": "Kendall", "last_name": "Roberts",
	}
	jsonBody, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	loginBody := map[string]string{"email": "test@example.com", "password": "securepassword123"}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"].(string)
}

// --- Registration Tests ---

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body := map[string]string{
		"email":      "test@example.com",
		"password":   "securepassword123",
		"first_name": "Kendall",
		"last_name":  "Roberts",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", response["email"])
	}
}

func TestRegister_BlockedWhenUserExists(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body := map[string]string{
		"email":      "first@example.com",
		"password":   "securepassword123",
		"first_name": "First",
		"last_name":  "User",
	}
	jsonBody, _ := json.Marshal(body)

	// First registration should succeed
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("First registration failed: %d", w.Code)
	}

	// Second registration should be blocked
	body["email"] = "second@example.com"
	jsonBody, _ = json.Marshal(body)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Login Tests ---

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Register
	regBody := map[string]string{
		"email": "test@example.com", "password": "securepassword123",
		"first_name": "Kendall", "last_name": "Roberts",
	}
	jsonBody, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Login
	loginBody := map[string]string{"email": "test@example.com", "password": "securepassword123"}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["token"] == nil {
		t.Error("Expected token in response")
	}
	if response["mfa_required"] != false {
		t.Error("Expected mfa_required to be false")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Register
	regBody := map[string]string{
		"email": "test@example.com", "password": "securepassword123",
		"first_name": "Kendall", "last_name": "Roberts",
	}
	jsonBody, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Login with wrong password
	loginBody := map[string]string{"email": "test@example.com", "password": "wrongpassword"}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// --- Token Refresh Test ---

func TestRefresh_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["token"] == nil {
		t.Error("Expected new token")
	}
}

// --- MFA Setup Test ---

func TestMFASetup_Success(t *testing.T) {
	setupTestDB(t)

	// Need crypto key for TOTP encryption
	testKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if err := crypto.SetKey(testKey); err != nil {
		t.Fatalf("Failed to set encryption key: %v", err)
	}

	router := setupRouter()
	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/mfa/setup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["qr_url"] == nil {
		t.Error("Expected qr_url in response")
	}
	if response["recovery_codes"] == nil {
		t.Error("Expected recovery_codes in response")
	}
	codes := response["recovery_codes"].([]interface{})
	if len(codes) != 8 {
		t.Errorf("Expected 8 recovery codes, got %d", len(codes))
	}
}
