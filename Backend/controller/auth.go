package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"webadmin/config"

	"github.com/gin-gonic/gin"
)

type tokenClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

func secretKey() []byte {
	secret := config.Get("jwt_secret")
	if secret == "" {
		secret = "your_jwt_secret_key"
	}
	return []byte(secret)
}

func GenerateToken(userID int, username, email string) (string, error) {
	claims := tokenClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Iat:      time.Now().Unix(),
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	payloadPart := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secretKey())
	_, _ = mac.Write([]byte(payloadPart))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return payloadPart + "." + signature, nil
}

func ParseToken(token string) (*tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}

	payloadPart := parts[0]
	signaturePart := parts[1]

	mac := hmac.New(sha256.New, secretKey())
	_, _ = mac.Write([]byte(payloadPart))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signaturePart), []byte(expectedSignature)) {
		return nil, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return nil, err
	}

	var claims tokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := ""
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
		if token == "" {
			token = strings.TrimSpace(c.GetHeader("token"))
		}
		if token == "" {
			c.JSON(401, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Next()
	}
}
