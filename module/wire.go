//go:build wireinject
// +build wireinject

package module

import (
	"okra_board2/repositories"
	"okra_board2/services"
	"okra_board2/controllers"
    "okra_board2/middlewares"
	"gorm.io/gorm"
	"github.com/google/wire"
)


func InitAdminController(db *gorm.DB) (c controllers.AdminController) {
    wire.Build(
        repositories.NewAdminRepositoryImpl,
        repositories.NewAuthRepositoryImpl,
        services.NewAdminServiceImpl,
        services.NewAuthServiceImpl,
        controllers.NewAdminControllerImpl, 
    )
    return
}

func InitAdminController2(db *gorm.DB) (c controllers.AdminController) {
    wire.Build(
        wire.Bind(
            repositories.NewAdminRepositoryImpl,
            services.NewAdminServiceImpl,
        ),
        wire.Bind(
            repositories.NewAuthRepositoryImpl,
            services.NewAuthServiceImpl,
        ),
        controllers.NewAdminControllerImpl,
    )
    return
}

func InitAuthMiddleware(db *gorm.DB) (m middlewares.AuthMiddleware) {
    wire.Build(
        repositories.NewAuthRepositoryImpl,
        services.NewAuthServiceImpl,
        middlewares.NewAuthMiddlewareImpl,
    )
    return
}
