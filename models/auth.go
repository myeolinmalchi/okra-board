package models

type AdminAuth struct {
    UUID            string      `gorm:"primaryKey;<-:create"`
    AdminID         string      `gorm:"<-:create"`
    AccessToken     string
    RefreshToken    string      `gorm:"<-:create"`
}
