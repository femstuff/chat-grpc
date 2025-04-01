package jwt

import (
	"errors"
	"time"

	"chat-grpc/Auth-service/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JWTService struct {
	Key           string
	TokenDuration time.Duration
	log           *zap.Logger
}

type Claims struct {
	UserID int64
	Role   string
	jwt.RegisteredClaims
}

func NewJWTService(key string, tokenDuration time.Duration, log *zap.Logger) *JWTService {
	return &JWTService{
		Key:           key,
		TokenDuration: tokenDuration,
		log:           log,
	}
}

func (j *JWTService) GenerateAccessToken(userID int64, role entity.Role) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role.StringRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Key))
}

func (j *JWTService) GenerateRefreshToken(userID int64, role entity.Role) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role.StringRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Key))
}

func (j *JWTService) VerifyAccessToken(tokenStr string) (*Claims, error) {
	j.log.Info(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Key), nil
	})
	if err != nil {
		j.log.Info("Error with parsing token", zap.String("token", tokenStr), zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("incorrect token")
	}

	return claims, nil
}

func (j *JWTService) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}
