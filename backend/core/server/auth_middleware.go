package server

import (
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const ContextParamUser = "CONTEXT_USER"

type AuthMiddleware struct {
	jwtService security.JWTService
}

func NewAuthMiddleware(jwtService security.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) RequireAuth(c *gin.Context) {

	//Get the  bearer Token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization token is required"})
		c.Abort()
		return
	}

	extractedToken := strings.Split(tokenString, "Bearer ")

	if len(extractedToken) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect format of authorization token"})
		c.Abort()
		return
	}

	claims, err := m.jwtService.ParseAndValidate(strings.TrimSpace(extractedToken[1]))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Token is not valid"})
		c.Abort()
		return
	}

	if time.Now().Unix() > claims.ExpiresAt {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.Abort()
		return
	}

	c.Set(ContextParamUser, claims)
	c.Next()
}

func GetContextClaims(c *gin.Context) (*security.JWTClaim, error) {
	s, ok := c.Get(ContextParamUser)
	if !ok {
		return nil, utils.ErrorPermission.New("wrong session")
	}
	claims, ok := s.(*security.JWTClaim)
	if !ok {
		return nil, utils.ErrorPermission.New("broken session")
	}
	return claims, nil
}
