package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kendall1978/project_salary/backend/internal/crypto"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret []byte

func SetJWTSecret(secret string) {
	JWTSecret = []byte(secret)
}

// --- Request types ---

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type MFAVerifyRequest struct {
	ChallengeToken string `json:"challenge_token" binding:"required"`
	Code           string `json:"code" binding:"required"`
}

type MFARecoverRequest struct {
	ChallengeToken string `json:"challenge_token" binding:"required"`
	RecoveryCode   string `json:"recovery_code" binding:"required"`
}

// --- JWT ---

func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// --- Handlers ---

func AuthStatus(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"registration_open": count == 0})
}

func Register(c *gin.Context) {
	// Single-user system: block registration if any user exists
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is disabled"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if MFA is enabled
	var mfa models.UserMFA
	mfaEnabled := database.DB.Where("user_id = ? AND is_enabled = ?", user.ID, true).First(&mfa).Error == nil

	if mfaEnabled {
		// Generate an opaque challenge token instead of exposing user_id
		challengeBytes := make([]byte, 32)
		rand.Read(challengeBytes)
		challengeToken := hex.EncodeToString(challengeBytes)

		// Clean up any expired challenges for this user
		database.DB.Where("user_id = ? OR expires_at < ?", user.ID, time.Now()).Delete(&models.MFAChallenge{})

		// Store challenge with 5-minute TTL
		database.DB.Create(&models.MFAChallenge{
			UserID:         user.ID,
			ChallengeToken: challengeToken,
			ExpiresAt:      time.Now().Add(5 * time.Minute),
		})

		c.JSON(http.StatusOK, gin.H{
			"mfa_required":    true,
			"challenge_token": challengeToken,
		})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"mfa_required": false,
	})
}

func MFASetup(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SalaryDashboard",
		AccountName: user.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP key"})
		return
	}

	// Encrypt the TOTP secret before storing
	encryptedSecret, err := crypto.Encrypt(key.Secret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt TOTP secret"})
		return
	}

	// Save MFA record with encrypted secret
	mfa := models.UserMFA{
		UserID:     userID,
		TOTPSecret: encryptedSecret,
		IsEnabled:  true,
	}
	database.DB.Where("user_id = ?", userID).Delete(&models.UserMFA{})
	database.DB.Create(&mfa)

	// Generate recovery codes
	recoveryCodes := make([]string, 8)
	database.DB.Where("user_id = ?", userID).Delete(&models.RecoveryCode{})
	for i := 0; i < 8; i++ {
		code := generateRecoveryCode()
		recoveryCodes[i] = code
		hash, _ := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		database.DB.Create(&models.RecoveryCode{
			UserID:   userID,
			CodeHash: string(hash),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_url":         key.URL(),
		"secret":         key.Secret(),
		"recovery_codes": recoveryCodes,
	})
}

func generateRecoveryCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

const maxMFAAttempts = 5

// lookupChallenge validates the challenge token, checks expiry and rate limits.
func lookupChallenge(c *gin.Context, challengeToken string) (*models.MFAChallenge, uint, bool) {
	var challenge models.MFAChallenge
	if err := database.DB.Where("challenge_token = ?", challengeToken).First(&challenge).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired challenge"})
		return nil, 0, false
	}

	if time.Now().After(challenge.ExpiresAt) {
		database.DB.Delete(&challenge)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Challenge expired, please login again"})
		return nil, 0, false
	}

	if challenge.FailedAttempts >= maxMFAAttempts {
		database.DB.Delete(&challenge)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many failed attempts, please login again"})
		return nil, 0, false
	}

	return &challenge, challenge.UserID, true
}

func MFAVerify(c *gin.Context) {
	var req MFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	challenge, userID, ok := lookupChallenge(c, req.ChallengeToken)
	if !ok {
		return
	}

	var mfa models.UserMFA
	if err := database.DB.Where("user_id = ? AND is_enabled = ?", userID, true).First(&mfa).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not enabled"})
		return
	}

	// Decrypt the stored TOTP secret before validating
	decryptedSecret, err := crypto.Decrypt(mfa.TOTPSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt TOTP secret"})
		return
	}

	if !totp.Validate(req.Code, decryptedSecret) {
		database.DB.Model(challenge).Update("failed_attempts", challenge.FailedAttempts+1)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	// Challenge used successfully — delete it
	database.DB.Delete(challenge)

	token, err := generateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func MFARecover(c *gin.Context) {
	var req MFARecoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	challenge, userID, ok := lookupChallenge(c, req.ChallengeToken)
	if !ok {
		return
	}

	var codes []models.RecoveryCode
	database.DB.Where("user_id = ? AND used_at IS NULL", userID).Find(&codes)

	for _, code := range codes {
		if bcrypt.CompareHashAndPassword([]byte(code.CodeHash), []byte(req.RecoveryCode)) == nil {
			now := time.Now()
			database.DB.Model(&code).Update("used_at", &now)

			// Challenge used successfully — delete it
			database.DB.Delete(challenge)

			token, err := generateToken(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
	}

	database.DB.Model(challenge).Update("failed_attempts", challenge.FailedAttempts+1)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid recovery code"})
}

func RefreshToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	token, err := generateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
