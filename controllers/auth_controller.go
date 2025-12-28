package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser mendaftarkan user baru + hashing password
func RegisterUser(c *gin.Context) {
	var input struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		FullName     string `json:"full_name"`
		Role         string `json:"role"`
		DepartmentID *uint  `json:"department_id"`
		Email        string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		FullName:     input.FullName,
		Role:         input.Role,
		DepartmentID: input.DepartmentID,
		Email:        input.Email,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username/Email likely exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration success", "user": user.Username})
}

// LoginUser mengecek password
func LoginUser(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Di sini Anda bisa generate JWT Token jika perlu. Untuk sekarang return sukses saja.
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "role": user.Role, "user_id": user.ID})
}