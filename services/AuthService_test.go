package services_test

import (
	"okra_board2/config"
	"okra_board2/models"
	"okra_board2/repositories"
	"okra_board2/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {
    db, err := config.InitDBConnection()
    if err != nil { assert.Error(t, err) }
    authRepo := repositories.NewAuthRepositoryImpl(db)
    adminRepo := repositories.NewAdminRepositoryImpl(db)
    adminService := services.NewAdminServiceImpl(adminRepo)
    authService := services.NewAuthServiceImpl(authRepo, adminService)

    admin := models.Admin {
        ID: "administrator11",
        Password: "@@Test123456",
        Email: "test123456@gmail.com",
        Name: "강민석",
        Phone: "010-4321-4321",
    }
    adminService.Register(&admin)

    auth, err := authService.CreateTokenPair("administrator11")
    if err != nil { assert.Error(t, err) }
    atClaims, err := authService.VerifyAccessToken(auth.AccessToken)
    if err != nil { assert.Error(t, err) }
    rtClaims, err := authService.VerifyRefreshToken(auth.RefreshToken)
    if err != nil { assert.Error(t, err) }

    uuid, err := authService.VerifyTokenPair(auth.AccessToken, auth.RefreshToken)
    assert.Equal(t, "administrator11", atClaims["id"].(string))
    assert.Equal(t, "강민석", atClaims["name"].(string))
    assert.Equal(t, "administrator11", rtClaims["id"].(string))
    assert.Equal(t, "강민석", rtClaims["name"].(string))
    assert.Equal(t, uuid, auth.UUID)

    time.Sleep(time.Second * 1)

    at, err := authService.CreateAccessToken(auth.UUID, auth.AdminID)
    atClaims, err = authService.VerifyAccessToken(at)
    assert.Equal(t, err, nil)
    assert.Equal(t, "administrator11", atClaims["id"].(string))
    assert.Equal(t, "강민석", atClaims["name"].(string))
    assert.Equal(t, auth.UUID, atClaims["uuid"].(string))

    assert.NotEqual(t, at, auth.AccessToken)

    _, err = authService.VerifyTokenPair(at, auth.RefreshToken)
    assert.Equal(t, err, nil)
    _, err = authService.VerifyTokenPair(auth.AccessToken, auth.RefreshToken)
    assert.EqualError(t, err, "Invalid Token Pair.")
    

    adminService.DeleteAdmin(admin.ID)

}
