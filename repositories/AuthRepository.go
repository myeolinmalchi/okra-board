package repositories

import (
	"okra_board2/models"
	"gorm.io/gorm"
)


type AuthRepository interface {

    // Select AdminAuth and returns with error
    GetAdminAuth(uuid string)                       (auth *models.AdminAuth, err error)

    // Insert AdminAuth and returns error
    InsertAdminAuth(adminAuth *models.AdminAuth)    (err error)

    // Delete AdminAuth and returns error
    DeleteAdminAuth(uuid string)                    (err error)

    // Update AdminAuth(Only Access Token) and returns error
    UpdateAccessToken(uuid, at string)              (err error)

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
    return }
