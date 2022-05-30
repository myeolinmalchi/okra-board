package controllers

import (
	"okra_board2/models"
	"okra_board2/services"
	"github.com/gin-gonic/gin"
)

type AdminController interface {
    Register(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
}

type AdminControllerImpl struct {
    adminService    services.AdminService
}

func NewAdminControllerImpl(
    adminService services.AdminService, 
) AdminController{
    return &AdminControllerImpl{ 
        adminService: adminService,
    }
}

func (a *AdminControllerImpl) Register(c *gin.Context) {
    requestBody := &models.Admin{}
    err := c.ShouldBind(requestBody)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    ok, result := a.adminService.Register(requestBody)
    if ok {
        c.Status(200)
    } else if result == nil {
        c.Status(400)
    } else {
        c.IndentedJSON(422, result)
    }
}

func (a *AdminControllerImpl) Update(c *gin.Context) {

    requestBody := &models.Admin{}
    err := c.ShouldBind(requestBody)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    ok, result := a.adminService.Update(requestBody)
    if ok {
        c.Status(200)
    } else if result == nil {
        c.Status(400)
    } else {
        c.IndentedJSON(422, result)
    }
}

func (a *AdminControllerImpl) Delete(c *gin.Context) { 
    
}
