package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	adminJwtSecret []byte
	appJwtSecret   []byte
)

func init() {
	// 从环境变量中读取 JWT 密钥，如果未设置则使用默认值
	adminSecret := os.Getenv("JWT_ADMIN_SECRET")
	if adminSecret == "" {
		adminSecret = "admin-secret-key"
	}
	adminJwtSecret = []byte(adminSecret)

	appSecret := os.Getenv("JWT_APP_SECRET")
	if appSecret == "" {
		appSecret = "app-secret-key"
	}
	appJwtSecret = []byte(appSecret)
}

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
