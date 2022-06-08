package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var Logger *log.Logger

func InitDBConnection() (*gorm.DB, error){
    dsn := "root:382274@tcp(localhost:3306)/board?parseTime=true"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            ti, _ := time.LoadLocation("Asia/Seoul")
            return time.Now().In(ti)
        },
    })
    return db, err
}

func InitLogger() (*os.File, error) {
    startTime := time.Now().Format("2006-01-02T15_04_05")
    fileName := "log/log-" + startTime
    if file, err := os.Create(strings.TrimSpace(fileName)); err != nil {
        return nil, err
    } else {
        Logger = log.New(file, "INFO: ", log.LstdFlags)
        return file, nil
    }
}

type Config struct {
    WhiteList       []string    `json:"whitelist"`
    AccessSecret    string      `json:"access_secret"`
    RefreshSecret   string      `json:"refresh_secret"`
    Domain          string      `json:"domain"`
}

func LoadConfig() (*Config, error){
    file, err := os.Open("config.json")
    defer file.Close()
    config := &Config{}
    jsonParser := json.NewDecoder(file)
    jsonParser.Decode(config)
    return config, err
}
