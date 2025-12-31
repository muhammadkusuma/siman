package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile (Sudah ada, biarkan seperti ini)
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var user models.User
	if err := models.DB.Preload("Department").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetAllUsers (Sudah ada)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	models.DB.Preload("Department").Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// --- TAMBAHAN BARU DI BAWAH SINI ---

// Input Struct untuk validasi update profil
type UpdateProfileInput struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UpdateProfile: Mengubah biodata user
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User

	// 1. Cari User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 2. Validasi Input
	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Update Data
	user.FullName = input.FullName
	user.Username = input.Username
	user.Email = input.Email

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "Profil berhasil diperbarui"})
}

// Input Struct untuk ganti password
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// ChangePassword: Mengganti password user
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User

	// 1. Cari User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 2. Validasi Input JSON
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Cek Password Lama (Gunakan PasswordHash, BUKAN Password)
	//
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama salah!"})
		return
	}

	// 4. Hash Password Baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password baru"})
		return
	}

	// 5. Simpan Password Baru (Gunakan PasswordHash, BUKAN Password)
	user.PasswordHash = string(hashedPassword)

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
