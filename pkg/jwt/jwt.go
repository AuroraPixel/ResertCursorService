package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	adminJwtSecret = []byte("admin-secret-key") // 管理员token密钥
	appJwtSecret   = []byte("app-secret-key")   // app token密钥
)

type Claims struct {
	AdminID uint `json:"adminId"`
	jwt.RegisteredClaims
}

type AppClaims struct {
	CodeID uint `json:"codeId"`
	jwt.RegisteredClaims
}

func GenerateToken(adminID uint) (string, error) {
	claims := Claims{
		AdminID: adminID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(adminJwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return adminJwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GenerateAppToken 生成 app 端使用的 token
func GenerateAppToken(codeID uint) (string, error) {
	claims := AppClaims{
		CodeID: codeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(appJwtSecret)
}

// ParseAppToken 解析 app 端的 token
func ParseAppToken(tokenString string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return appJwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
