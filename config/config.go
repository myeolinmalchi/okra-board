package config

import (
	"encoding/json"
	"fmt"
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

type Config struct {
    WhiteList       []string    `json:"whitelist"`
    AccessSecret    string      `json:"access_secret"`
    RefreshSecret   string      `json:"refresh_secret"`
    Domain          string      `json:"domain"`
    DefaultThumbnail string     `json:"default_thumbnail"`
    DB              DBConfig    `json:"db"`
    Log             LogConfig   `json:"log"`
    AWS             AWSConfig   `json:"aws"`
}

type DBConfig struct {
    User        string          `json:"user"`
    Password    string          `json:"password"`
    Host        string          `json:"host"`
    Port        int             `json:"port"`
    Database    string          `json:"database"`
}

func (c *DBConfig) ToString() string {
    return fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?parseTime=true",
        c.User,
        c.Password,
        c.Host,
        c.Port,
        c.Database,
    )
}

type LogConfig struct {
    Path        string          `json:"path"`
    TimeFormat  string          `json:"time_format"`
    Prefix      string          `json:"prefix"`
}

type AWSConfig struct {
    AccessKeyID string          `json:"access_key_id"`
    SecretKey   string          `json:"secret_key"`
    Region      string          `json:"region"`
    Bucket      string          `json:"bucket"`
    Domain      string          `json:"domain"`
}

func LoadConfig() (*Config, error){
    file, err := os.Open("config.json")
    defer file.Close()
    config := &Config{}
    jsonParser := json.NewDecoder(file)
    jsonParser.Decode(config)
    return config, err
}

func LoadConfigTest() (*Config, error){
    file, err := os.Open("../config.json")
    defer file.Close()
    config := &Config{}
    jsonParser := json.NewDecoder(file)
    jsonParser.Decode(config)
    return config, err
}


func InitDBConnection(conf *Config) (*gorm.DB, error){
    dsn := conf.DB.ToString()
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            ti, _ := time.LoadLocation("Asia/Seoul")
            return time.Now().In(ti)
        },
    })
    return db, err
}

func InitLogger(conf *Config) (*os.File, error) {
    startTime := time.Now().Format(conf.Log.TimeFormat)
    fileName := conf.Log.Path + "/" + conf.Log.Prefix + "-" + startTime
    if file, err := os.Create(strings.TrimSpace(fileName)); err != nil {
        return nil, err
    } else {
        Logger = log.New(file, "INFO: ", log.LstdFlags)
        return file, nil
    }
}
