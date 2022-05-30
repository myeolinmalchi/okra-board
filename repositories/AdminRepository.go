package repositories

import (
	"okra_board2/models"
	"gorm.io/gorm"
    "okra_board2/utils/encryption"
)

type AdminRepository interface {
    
    // 해당 열의 값에 대응되는 관리자 계정이 존재하는지 체크한다.
    CheckAdminExists(column, value string)  bool

    // Select Admin Account and returns with error
    GetAdmin(id string)                     (admin *models.Admin, err error)

    // Insert Admin Account and returns error
    InsertAdmin(*models.Admin)              error

    // Update Admin Account and returns error
    UpdateAdmin(*models.Admin)              error

    // Delete Admin Account and returns error
    DeleteAdmin(string)                     error
}

type AdminRepositoryImpl struct {
    db *gorm.DB
}

func NewAdminRepositoryImpl(db *gorm.DB) AdminRepository {
    return &AdminRepositoryImpl{ db: db }
}

func (rep *AdminRepositoryImpl) CheckAdminExists(column, value string) (exists bool) {
    rep.db.Model(&models.Admin{}).
        Select("count(*) > 0").
        Where(column + " = ?", value).
        Find(&exists)
    return
}

func (rep *AdminRepositoryImpl) GetAdmin(id string) (admin *models.Admin, err error) {
    admin = &models.Admin{}
    err = rep.db.Table("admins").First(admin, "id = ?", id).Error
    return
}

func (rep *AdminRepositoryImpl) InsertAdmin(admin *models.Admin) (err error) {
    admin.Password = encryption.EncryptSHA256(admin.Password)
    err = rep.db.Create(admin).Error
    return
}

func (rep *AdminRepositoryImpl) UpdateAdmin(admin *models.Admin) (err error) {
    admin.Password = encryption.EncryptSHA256(admin.Password)
    err = rep.db.UpdateColumns(admin).Error
    return
}

func (rep *AdminRepositoryImpl) DeleteAdmin(id string) (err error) {
    err = rep.db.Delete(&models.Admin{}, "id", id).Error
    return nil
}

