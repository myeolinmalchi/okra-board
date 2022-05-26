package repositories

import (
	"okra_board2/models"
	"gorm.io/gorm"
)


type AuthRepository interface {

    // Insert AdminAuth
    InsertAdminAuth(adminAuth *models.AdminAuth) error 

    // Select AdminAuth
    GetAdminAuth(uuid string) (*models.AdminAuth, error)

    // Delete AdminAuth
    DeleteAdminAuth(uuid string) error

    // Update AdminAuth(Only Access Token)
    UpdateAccessToken(uuid, at string) error

}

type AuthRepositoryImpl struct {
    db *gorm.DB
}

// AuthRepositoryImpl 객체를 생성한다.
func NewAuthRepositoryImpl(db *gorm.DB) AuthRepository {
    return &AuthRepositoryImpl{ db: db }
}

func (rep *AuthRepositoryImpl) InsertAdminAuth(adminAuth *models.AdminAuth) (err error) {
    err = rep.db.Create(adminAuth).Error
    return
}

func (rep *AuthRepositoryImpl) GetAdminAuth(uuid string) (adminAuth *models.AdminAuth, err error) {
    err = rep.db.First(&adminAuth, "uuid = ?", uuid).Error
    return
}

func (rep *AuthRepositoryImpl) DeleteAdminAuth(uuid string) (err error) {
    err = rep.db.Delete(&models.AdminAuth{}, "uuid = ?", uuid).Error
    return
}

func (rep *AuthRepositoryImpl) UpdateAccessToken(uuid, at string) (err error) {
    err = rep.db.Table("admin_auths").
        Where("uuid = ?", uuid).
        Update("access_token", at).
        Error
    return
}
