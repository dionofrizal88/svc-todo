package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

type UserLoginClaims struct {
    jwt.StandardClaims
    Username string `json:"username"`
    Role     string `json:"role"`
}

func (manager *JWTManager) Generate(user *UserLogin) (string, error) {
    claims := UserLoginClaims{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
        },
        Username: user.Username,
        Role:     user.Role,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(accessToken string) (*UserLoginClaims, error) {
    token, err := jwt.ParseWithClaims(
        accessToken,
        &UserLoginClaims{},
        func(token *jwt.Token) (interface{}, error) {
            _, ok := token.Method.(*jwt.SigningMethodHMAC)
            if !ok {
                return nil, fmt.Errorf("unexpected token signing method")
            }

            return []byte(manager.secretKey), nil
        },
    )

    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    claims, ok := token.Claims.(*UserLoginClaims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }

    return claims, nil
}