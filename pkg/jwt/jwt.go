package jwt

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat JWT baru dengan userID sebagai subject
func GenerateToken(userID uint, secret string, duration time.Duration) (string, error) {
	log.Printf("[JWT] Generating token for userID: %d, secret length: %d", userID, len(secret))

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Printf("[JWT] Error generating token: %v", err)
		return "", err
	}

	log.Printf("[JWT] Token generated successfully for userID: %d", userID)
	return tokenString, nil
}

// AuthMiddleware memverifikasi JWT dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Cek Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("[JWT] Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		// 2. Extract token dari Bearer
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			// Tidak ada prefix "Bearer ", format salah
			log.Println("[JWT] Invalid token format: missing 'Bearer ' prefix")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format, use: Bearer <token>"})
			c.Abort()
			return
		}

		if tokenStr == "" {
			log.Println("[JWT] Empty token after Bearer prefix")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}

		// 3. Get secret dari environment
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("[JWT] CRITICAL: JWT_SECRET not set in environment")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
			c.Abort()
			return
		}

		log.Printf("[JWT] Verifying token (secret length: %d, token length: %d)", len(secret), len(tokenStr))
		log.Printf("[JWT] Token received: %s...%s", tokenStr[:20], tokenStr[len(tokenStr)-20:])

		// 4. Parse dan validasi token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// Validasi signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("[JWT] Unexpected signing method: %v", t.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		// 5. Handle parsing errors
		if err != nil {
			log.Printf("[JWT] Token parse error: %v", err)

			// Berikan pesan error yang lebih spesifik
			var errorMsg string
			switch {
			case err.Error() == "token is expired":
				errorMsg = "token has expired, please login again"
			case strings.Contains(err.Error(), "signature is invalid"):
				errorMsg = "invalid token signature"
			default:
				errorMsg = "invalid token"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  errorMsg,
				"detail": err.Error(),
			})
			c.Abort()
			return
		}

		// 6. Cek validitas token
		if !token.Valid {
			log.Println("[JWT] Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			c.Abort()
			return
		}

		// 7. Extract userID dari claims
		sub, ok := claims["sub"].(float64)
		if !ok {
			log.Printf("[JWT] Invalid sub claim type: %T", claims["sub"])
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userID := uint(sub)
		log.Printf("[JWT] ✓ Token validated successfully for userID: %d", userID)

		// 8. Set userID di context
		c.Set("userID", userID)

		// 9. Continue ke handler berikutnya
		c.Next()
	}
}

// GetUserIDFromContext helper function untuk mengambil userID dari context
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("userID not found in context")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid userID type in context")
	}

	return id, nil
}
