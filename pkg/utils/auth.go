package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	Key string
}

type Claims struct {
	jwt.StandardClaims
	OrgId int64 `json:"user_id"`
}

func NewTokenService(key string) *TokenService {
	return &TokenService{key}
}

func (ts *TokenService) CreateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&Claims{
			jwt.StandardClaims{},
			id,
		})
	return token.SignedString([]byte(ts.Key))
}

func (ts *TokenService) ParseToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(ts.Key), nil
		})
	if err != nil{
		return 0, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	return claims.OrgId, nil
}