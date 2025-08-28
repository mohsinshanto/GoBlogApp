package controllers

import (
	"BlogApp/config"
	"BlogApp/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(c *gin.Context) {
	// Define input struct with only required fields
	type RegisterInput struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input", "error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "Username or Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to process password"})
		return
	}

	// Create new user
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Registration successful"})
}

func Login(c *gin.Context) {
	type LoginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input", "error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret) // jwtSecret must be []byte
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Token creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "Login successful",
		"token": tokenString,
	})
}
func GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	userID := uint(rawID.(float64))

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
		return
	}

	// Return user info (omit password)
	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"profileImage": user.ProfileImage,
	})
}
func UpdateProfile(c *gin.Context) {
	// Get user ID from context
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	userID := uint(rawID.(float64))

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
		return
	}

	// Parse form-data for image
	username := c.PostForm("username")
	email := c.PostForm("email")

	// Update profile image if uploaded
	file, err := c.FormFile("profileImage")
	if err == nil {
		// Save the file locally (or to cloud)
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to upload image"})
			return
		}
		user.ProfileImage = "/" + path
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"profileImage": user.ProfileImage,
	})
}
