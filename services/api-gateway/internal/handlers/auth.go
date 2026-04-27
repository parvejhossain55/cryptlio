package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=30"`
	Bio      string `json:"bio" binding:"omitempty,max=200"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	KYCLevel int    `json:"kyc_level"`
}

// In-memory user store for development (use database in production)
var devUsers = make(map[string]User)

func init() {
	// Add a test user
	devUsers["user@example.com"] = User{
		ID:       "user_001",
		Email:    "user@example.com",
		Username: "testuser",
		KYCLevel: 2,
	}
}

func RegisterHandler(jwtSecret string, jwtExpiry time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"details": err.Error(),
			})
			return
		}

		// Check if user exists
		if _, exists := devUsers[req.Email]; exists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
			return
		}

		// Create user (simplified - no password hashing for dev)
		user := User{
			ID:       "user_" + time.Now().Format("20060102150405"),
			Email:    req.Email,
			Username: req.Username,
			KYCLevel: 1,
		}
		devUsers[req.Email] = user

		// Generate JWT
		token := generateJWT(user, jwtSecret, jwtExpiry)

		c.JSON(http.StatusCreated, AuthResponse{
			Token: token,
			User:  user,
		})
	}
}

func LoginHandler(jwtSecret string, jwtExpiry time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"details": err.Error(),
			})
			return
		}

		// Find user
		user, exists := devUsers[req.Email]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// TODO: Verify password hash (accept any password for now in dev)
		// In production: bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

		// Generate JWT
		token := generateJWT(user, jwtSecret, jwtExpiry)

		c.JSON(http.StatusOK, AuthResponse{
			Token: token,
			User:  user,
		})
	}
}

// GetUserHandler returns current user profile (dev mock)
func GetUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_id from context (set by auth middleware)
		userID, _ := c.Get("user_id")

		// Find user by ID
		for _, user := range devUsers {
			if user.ID == userID {
				c.JSON(http.StatusOK, gin.H{
					"user": user,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	}
}

// UpdateUserHandler updates user profile (dev mock)
func UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"details": err.Error(),
			})
			return
		}

		// Find and update user
		for email, user := range devUsers {
			if user.ID == userID {
				if req.Username != "" {
					user.Username = req.Username
				}
				devUsers[email] = user

				c.JSON(http.StatusOK, gin.H{
					"user": user,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
	}
}

func generateJWT(user User, secret string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"kyc_level":  user.KYCLevel,
		"token_type": "access",
		"exp":        time.Now().Add(expiry).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}
