package tokens

import (
	"context"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTService struct{
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{AccessSecret: []byte(accessSecret), RefreshSecret: []byte(refreshSecret), AccessTTL: accessTTL, RefreshTTL: refreshTTL}
}

type claims struct{
	Type string `json:"type"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTService) GenerateAccessToken(ctx context.Context, userID, role string) (string, error) {
	c := claims{Type: "access", Role: role, RegisteredClaims: jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(j.AccessSecret)
}

func (j *JWTService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	c := claims{Type: "refresh", RegisteredClaims: jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.RefreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(j.RefreshSecret)
}

func (j *JWTService) ValidateAccessToken(ctx context.Context, token string) (userID string, role string, err error) {
	parsed, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) { return j.AccessSecret, nil })
	if err != nil { return "", "", err }
	c, ok := parsed.Claims.(*claims)
	if !ok || !parsed.Valid || c.Type != "access" { return "", "", errors.New("invalid token") }
	return c.Subject, c.Role, nil
}

func (j *JWTService) ValidateRefreshToken(ctx context.Context, token string) (userID string, err error) {
	parsed, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) { return j.RefreshSecret, nil })
	if err != nil { return "", err }
	c, ok := parsed.Claims.(*claims)
	if !ok || !parsed.Valid || c.Type != "refresh" { return "", errors.New("invalid token") }
	return c.Subject, nil
}
