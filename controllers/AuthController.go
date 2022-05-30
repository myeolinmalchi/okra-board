package controllers

import (
    "okra_board2/services"
    "okra_board2/models"
    "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
    "strings"
)

type AuthController interface {
    Auth(c *gin.Context)
    Login(c *gin.Context)
    ReissueAccessToken(c *gin.Context)
}

type AuthControllerImpl struct {
    authService services.AuthService
    adminService services.AdminService
}

func NewAuthControllerImpl(
    authService services.AuthService,
    adminService services.AdminService,
) AuthController {
    return &AuthControllerImpl{ 
        authService: authService,
        adminService: adminService,
    }
}

func (a *AuthControllerImpl) Auth(c *gin.Context) {
    authorization := c.Request.Header.Get("Authorization")
    tokenPair := strings.Split(authorization, " ")
    token := tokenPair[0] // access token
    if token == "" {
        c.JSON(401, gin.H {
            "status": 401,
            "message": "access token is empty.",
        })
    } else if _, err := a.authService.VerifyAccessToken(token); err != nil {
        if v, _ := err.(*jwt.ValidationError); v.Errors == jwt.ValidationErrorExpired {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "access token is expired.",
            })
        } else {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "invalid access token.",
            })
        }
    }
}

func (a *AuthControllerImpl) Login(c *gin.Context) {
    requestBody := &models.Admin{}
    err := c.ShouldBind(requestBody)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }

    if a.adminService.Login(requestBody) {
        adminAuth, err := a.authService.CreateTokenPair(requestBody.ID)
        if err != nil {
            c.JSON(400, err.Error())
            return
        }
        tokenPair := adminAuth.AccessToken + " " + adminAuth.RefreshToken
        c.Header("Authorization", tokenPair)
        c.Status(200)
    } else {
        c.Status(401)
    }
}

func (a *AuthControllerImpl) ReissueAccessToken(c *gin.Context) {
    authorization := c.Request.Header.Get("Authorization")
    tokenPair := strings.Split(authorization, " ")
    var accessToken, refreshToken string
    accessToken = tokenPair[0]
    if len(tokenPair) >= 2 {
        refreshToken = tokenPair[1]
    } else {
        c.JSON(401, gin.H {
            "status": 401,
            "message": "refresh token is empty",
        })
        return
    }

    claims, err := a.authService.VerifyRefreshToken(refreshToken)
    uuid, ok := claims["uuid"].(string)
    if !ok {
        c.JSON(401, gin.H {
            "status": 401,
            "message": "invalid refresh token",
        })
        return
    }
    id, ok := claims["id"].(string)
    if !ok {
        c.JSON(401, gin.H {
            "status": 401,
            "message": "invalid refresh token",
        })
        return
    }
    if err != nil {
        if err := a.authService.DeleteTokenPair(uuid); err != nil {
            c.JSON(401, gin.H {
                "status": 401,
                "message": err.Error(),
            })
        }
        if v, _ := err.(*jwt.ValidationError); v.Errors == jwt.ValidationErrorExpired {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "refresh token is expired.",
            })
        } else {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "invalid refresh token.",
            })
        }
    } else if _, err := a.authService.VerifyTokenPair(accessToken, refreshToken); err != nil {
        if err := a.authService.DeleteTokenPair(uuid); err != nil {
            c.JSON(401, gin.H {
                "status": 401,
                "message": err.Error(),
            })
        }
        c.JSON(401, gin.H {
            "status": 401,
            "message": err.Error(),
        })
    } else {
        newAccessToken, err := a.authService.CreateAccessToken(uuid, id)
        if err != nil {
            c.JSON(401, gin.H {
                "status": 401,
                "message": err.Error(),
            })
            return
        }
        newTokenPair := newAccessToken + " " + refreshToken
        c.Header("Authorization", newTokenPair)
        c.Status(200)
    }
}
