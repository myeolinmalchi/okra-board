package services_test

import (
	"okra_board2/config"
	"okra_board2/models"
	"okra_board2/repositories"
	"okra_board2/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminService(t *testing.T) {
    conf, err := config.LoadConfigTest()
    if err != nil { assert.Error(t, err) }

    db, err := config.InitDBConnection(conf)
    if err != nil { assert.Error(t, err) }
    adminRepo := repositories.NewAdminRepositoryImpl(db)
    s := services.NewAdminServiceImpl(adminRepo)

    var nilStr *string

    var admin models.Admin
    // register case 1: incorrect id
    admin = models.Admin {ID: "test"}
    ok, result := s.Register(&admin)
    assert.Equal(t, ok, false)
    assert.NotEqual(t, result.ID, nilStr)

    // register case 2: correct id
    admin = models.Admin {ID: "test12345"}
    ok, result = s.Register(&admin)
    assert.Equal(t, ok, false)
    assert.Equal(t, result.ID, nilStr)

    // register case 3: incorrect pw
    admin = models.Admin {Password: "test"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.NotEqual(t, result.Password, nilStr)

    // register case 4: correct pw
    admin = models.Admin {Password: "@@Test123456"}
    ok, result = s.Register(&admin)
    assert.Equal(t, ok, false)
    assert.Equal(t, result.Password, nilStr)

    // register case 5: incorrect email
    admin = models.Admin {Email: "test"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.NotEqual(t, result.Email, nilStr)

    // register case 6: correct email
    admin = models.Admin {Email: "test1234@gmail.com"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.Equal(t, result.Email, nilStr)

    // register case 7: incorrect name
    admin = models.Admin {Name: "ㄱㄴㄷ"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.NotEqual(t, result.Name, nilStr)

    // register case 8: correct name
    admin = models.Admin {Name: "강민석"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.Equal(t, result.Name, nilStr)

    // register case 9: incorrect phone
    admin = models.Admin {Phone: "test"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.NotEqual(t, result.Phone, nilStr)

    // register case 10: correct phone
    admin = models.Admin {Phone: "010-1234-0012"}
    ok, result = s.Register(&admin) 
    assert.Equal(t, ok, false)
    assert.Equal(t, result.Phone, nilStr)

    // register case 11: all correct
    var nilResult *models.AdminValidationResult
    nilResult = nil
    admin = models.Admin {
        ID: "administrator11",
        Password: "@@Test123456",
        Email: "test123456@gmail.com",
        Name: "강민석",
        Phone: "010-4321-4321",
    }
    ok, result = s.Register(&admin)
    assert.Equal(t, ok, true)
    assert.Equal(t, result, nilResult)

    // update
    admin = models.Admin {
        ID: "administrator11",
        Password: "@@Test1234567",
        Email: "test1234567@gmail.com",
        Name: "강민석",
        Phone: "010-4321-4321",
    }
    ok, result = s.Update(&admin)
    assert.Equal(t, ok, true)
    assert.Equal(t, result, nilResult)

    // select
    getAdmin, err := s.GetAdmin("administrator11")
    if err != nil { assert.Error(t, err) }
    assert.Equal(t, "test1234567@gmail.com", getAdmin.Email)

    // login case 1: incorrect
    loginInfo := models.Admin {
        ID: "administrator11",
        Password: "@@Test123456",
    }
    ok = s.Login(&loginInfo)
    assert.Equal(t, ok, false)

    // login case 2: correct
    loginInfo = models.Admin {
        ID: "administrator11",
        Password: "@@Test1234567",
    }
    ok = s.Login(&loginInfo)
    assert.Equal(t, ok, true)

    // delete
    err = s.DeleteAdmin("administrator11")
    if err != nil { assert.Error(t, err) }
    
}
