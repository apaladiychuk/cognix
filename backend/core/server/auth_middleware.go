package server

import (
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"context"
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

	//testClaim := security.Identity{
	//	AccessToken:  "",
	//	RefreshToken: "",
	//	User: &model.User{
	//		ID:       uuid.MustParse("9b63b4ec-20dc-4b10-8597-5b8d9039403e"),
	//		TenantID: uuid.MustParse("c810016b-9506-4db2-ad40-a8e3aa517108"),
	//	},
	//}
	//c.Request = c.Request.WithContext(context.WithValue(
	//	c.Request.Context(), ContextParamUser, &testClaim))
	//c.Next()
	//
	//return

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

	c.Request = c.Request.WithContext(context.WithValue(
		c.Request.Context(), ContextParamUser, claims))
	c.Next()
}

func GetContextIdentity(c *gin.Context) (*security.Identity, error) {
	claims, ok := c.Request.Context().Value(ContextParamUser).(*security.Identity)
	if !ok {
		return nil, utils.ErrorPermission.New("broken session")
	}
	return claims, nil
}
