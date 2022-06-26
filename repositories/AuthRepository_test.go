package repositories_test

import (
	"okra_board2/config"
	"okra_board2/models"
	"okra_board2/repositories"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthCRUD(t *testing.T){ 
    conf, err := config.LoadConfigTest()
    if err != nil { assert.Error(t, err) }
    db, err := config.InitDBConnection(conf)
    if err != nil { assert.Error(t, err) }

    sqlDB, err := db.DB()
    defer sqlDB.Close()

    r := repositories.NewAuthRepositoryImpl(db)
    auths := make([]models.AdminAuth, 5)
    for i := 0; i < 5; i++ {
        index := strconv.Itoa(i+1)
        auths[i] = models.AdminAuth {
            UUID: "uuid" + index,
            AdminID: "okraseoul",
            AccessToken: "access token " + index,
            RefreshToken: "refresh token " + index,
        }
    }

    // insert
    for i := 0; i < 5; i++ {
        if err := r.InsertAdminAuth(&auths[i]); err != nil {
            assert.Error(t, err)
        }
    }

    // update
    err = r.UpdateAccessToken("uuid1", "updated access token")
    if err != nil { assert.Error(t, err) }

    // select
    auth, err := r.GetAdminAuth("uuid1")
    if err != nil { assert.Error(t, err) }
    assert.Equal(t, "updated access token", auth.AccessToken)

    // delete
    for i := 0; i < 5; i++ {
        if err := r.DeleteAdminAuth(auths[i].UUID); err != nil {
            assert.Error(t, err)
        }
    }
}
