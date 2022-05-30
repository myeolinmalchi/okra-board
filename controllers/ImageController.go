package controllers

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type ImageController interface {
    UploadImage(c *gin.Context)
}

type ImageControllerImpl struct {}

func NewImageControllerImpl() ImageController {
    return &ImageControllerImpl{}
}

func (i *ImageControllerImpl) UploadImage(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        log.Fatal(err)
        c.Status(400)
    }

    filename := uuid.NewString()
    if err := c.SaveUploadedFile(file, "./public/images/"+filename+".png"); err != nil {
        log.Fatal(err)
        c.Status(400)
    }

    domain := os.Getenv("DOMAIN")
    url := "http://"+domain+"/images/"+filename+".png"

    c.JSON(200, gin.H { "url": url })
}
