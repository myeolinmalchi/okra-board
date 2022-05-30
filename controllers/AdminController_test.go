package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"okra_board2/controllers"
	"okra_board2/models"
	"okra_board2/module"
	"testing"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var router *gin.Engine
var c controllers.AdminController
var db *gorm.DB

func init() {
    dsn := "root:382274@tcp(localhost:3306)/board_prototype?parseTime=true"
    db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            ti, _ := time.LoadLocation("Asia/Seoul")
            return time.Now().In(ti)
        },
    })
    c = module.InitAdminController(db)

    router = gin.Default()
    v1 := router.Group("/api/v1")
    {
        v1.POST("/admin/login", c.Login)
        v1.POST("/admin", c.Register)
        v1.PUT("/admin", c.Update)
        v1.POST("/admin/auth", c.ReissueAccessToken)
    }
}

func TestLogin(t *testing.T) {
    admin := models.Admin {
        ID: "minsuk4820",
        Password: "asd3254820",
    }
    jsonValue, _ := json.Marshal(admin)
    req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

   // tokenPair := w.Header().Get("Authorization")

}
