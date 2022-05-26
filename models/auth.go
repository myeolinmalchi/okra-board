package models

type AdminAuth struct {
    UUID            string      `gorm:"primaryKey"`
    AdminID         string
    AccessToken     string
    RefreshToken    string
}
