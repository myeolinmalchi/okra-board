package middlewares

import (
	"okra_board2/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthMiddleware interface {
    Auth(c *gin.Context)
}

type AuthMiddlewareImpl struct {
    authService services.AuthService
}

func NewAuthMiddlewareImpl(authService services.AuthService) AuthMiddleware {
    return &AuthMiddlewareImpl{ authService: authService }
}

type CustomGinContext struct {
    Context *gin.Context
}

func (c CustomGinContext) JSONWithStatus(status int, msg string) {
    c.Context.JSON(status, gin.H {
        "status": status,
        "message": msg,
    })
    c.Context.Abort()
}

func (m *AuthMiddlewareImpl) Auth(c *gin.Context) {
    context := CustomGinContext{ Context: c }
    authorization := c.Request.Header.Get("Authorization")
    tokenPair := strings.Split(authorization, " ")
    token := tokenPair[0] // access token
    if token == "" {
        context.JSONWithStatus(401, "access token is empty")
    } else if _, err := m.authService.VerifyAccessToken(token); err != nil {
        if v, _ := err.(*jwt.ValidationError); v.Errors == jwt.ValidationErrorExpired {
            context.JSONWithStatus(401, "access token is expired")
        } else {
            context.JSONWithStatus(401, "invalid access token")
        }
    }
}

