package jwt

import (
	"errors"
	"time"

	"chat-grpc/Auth-service/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Key           string
	tokenDuration time.Duration
}

type Claims struct {
	UserID int64
	Role   string
	jwt.RegisteredClaims
}

func NewJWTService(key string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		Key:           key,
		tokenDuration: tokenDuration,
	}
}

func (j *JWTService) GenerateToken(userID int64, role entity.Role) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role.StringRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Key))
}

func (j *JWTService) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("incorrect token")
	}

	return claims, nil
}
