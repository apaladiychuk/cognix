package security

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type (
	JWTClaim struct {
		jwt.StandardClaims
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserID       string `json:"user_id"`
		TenantID     string `json:"tenant_id"`
	}
	JWTService interface {
		Create(claim *JWTClaim) (string, error)
		ParseAndValidate(string) (*JWTClaim, error)
	}
	jwtService struct {
		jwtSecret      string `json:"jwt_secret"`
		jwtExpiredTime int    `json:"jwt_expired_time"`
	}
)

func NewJWTService(jwtSecret string, jwtExpiredTime int) JWTService {
	return &jwtService{jwtSecret: jwtSecret,
		jwtExpiredTime: jwtExpiredTime}
}

func (j *jwtService) Create(claim *JWTClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	claim.ExpiresAt = time.Now().Add(time.Duration(j.jwtExpiredTime)).Unix()
	tokenString, err := token.SignedString(j.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtService) ParseAndValidate(tokenString string) (*JWTClaim, error) {
	var claims JWTClaim
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
