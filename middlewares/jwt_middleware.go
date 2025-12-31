package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Kita hapus variabel global hardcoded
// var JwtSecret = []byte("RAHASIA_KAMPUS_SUPER_AMAN_2024")

// GetSecretKey adalah helper untuk mengambil kunci rahasia dari Environment Variable
func GetSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback value jika di .env kosong (Hanya untuk development, jangan dipakai di production!)
		return []byte("secret_default_kalau_lupa_set_env")
	}
	return []byte(secret)
}

// GenerateToken membuat token JWT
func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Gunakan fungsi GetSecretKey() di sini
	return token.SignedString(GetSecretKey())
}

// AuthMiddleware memvalidasi token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Gunakan fungsi GetSecretKey() di sini juga
			return GetSecretKey(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Simpan user_id ke context agar bisa dipakai di Controller
			if floatID, ok := claims["user_id"].(float64); ok {
				c.Set("userID", uint(floatID))
			}
			c.Set("role", claims["role"])
		}
		c.Next()
	}
}
