package main

import (
	"okra_board2/config"
	"okra_board2/module"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

    route := gin.Default()
    route.Use(cors.New(cors.Config {
        AllowAllOrigins:    true,
        AllowMethods:       []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:       []string{"Content-Type", "Authorization"},
        ExposeHeaders:      []string{"Authorization"},
        AllowCredentials:   true,
        MaxAge: 12 * time.Hour,
    }))

    db := config.NewDBConnection()

    authMiddleware := module.InitAuthMiddleware(db)

    adminController := module.InitAdminController(db)

    v1 := route.Group("/api/v1")
    {
        v1.POST("admin", authMiddleware.Auth, adminController.Register)
        v1.PUT("admin", authMiddleware.Auth, adminController.Update)
        v1.POST("admin/login", adminController.Login)
        v1.POST("admin/auth", adminController.ReissueAccessToken)
    }
    route.Run(":3000")
}
