//go:build wireinject
// +build wireinject

package module

import (
	"okra_board2/repositories"
	"okra_board2/services"
	"okra_board2/controllers"
	"gorm.io/gorm"
	"github.com/google/wire"
)


func InitAdminController(db *gorm.DB) (c controllers.AdminController) {
    wire.Build(
        repositories.NewAdminRepositoryImpl,
        services.NewAdminServiceImpl,
        controllers.NewAdminControllerImpl, 
    )
    return
}

func InitAuthController(db *gorm.DB) (a controllers.AuthController) {
    wire.Build(
        repositories.NewAuthRepositoryImpl,
        repositories.NewAdminRepositoryImpl,
        services.NewAdminServiceImpl,
        services.NewAuthServiceImpl,
        controllers.NewAuthControllerImpl,
    )
    return
}

func InitPostController(db *gorm.DB) (c controllers.PostController) {
    wire.Build( 
        repositories.NewPostRepositoryImpl,
        services.NewPostServiceImpl,
        controllers.NewPostControllerImpl,
    )
    return
}
