package repositories_test

import (
	"github.com/stretchr/testify/assert"
	"okra_board2/config"
	"okra_board2/models"
	"okra_board2/repositories"
	"strconv"
	"testing"
)

func TestAdminCRUD(t *testing.T) {
    db, err := config.InitDBConnection()
    if err != nil { assert.Error(t, err) }

    r := repositories.NewAdminRepositoryImpl(db)

    admins := make([]models.Admin, 5)
    for i := 0; i < 5; i++ {
        index := strconv.Itoa(i + 1)
        admins[i] = models.Admin {
            ID: "admin" + index,
            Password: "admin" + index,
            Name: "홍길동" + index,
            Email: "test"+index+"@gmail.com",
            Phone: "010-2274-382"+index,
        }
    }

    // insert 
    for i := 0; i < len(admins); i++ {
        if err := r.InsertAdmin(&admins[i]); err != nil {
            assert.Error(t, err)
        }
    }

    // update
    err = r.UpdateAdmin(&models.Admin {
        ID: "admin1",
        Password: "updated password",
        Name: "강민석",
        Email: "testemail@gmail.com",
        Phone: "010-2284-3820",
    })

    // check existence
    checkID := r.CheckAdminExists("id", "admin1")
    checkEmail := r.CheckAdminExists("email", "testemail@gmail.com")
    checkPhone := r.CheckAdminExists("phone", "010-2274-3820")

    assert.Equal(t, true, checkID)
    assert.Equal(t, true, checkEmail)
    assert.Equal(t, true, checkPhone)

    // select
    admin, err := r.GetAdmin("admin1")
    assert.Equal(t, "강민석", admin.Name)

    // delete
    for i := 0; i < len(admins); i++ {
        if err := r.DeleteAdmin(admins[i].ID); err != nil {
            assert.Error(t, err)
        }
    }
}
