package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfile mengambil data user yang sedang login berdasarkan Token JWT
func GetProfile(c *gin.Context) {
	// Ambil userID dari context (diset oleh Middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var user models.User
	// Preload Department untuk tahu dia dari unit mana
	if err := models.DB.Preload("Department").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetAllUsers (Biasanya untuk menu 'Manajemen User' bagi Admin)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	models.DB.Preload("Department").Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}