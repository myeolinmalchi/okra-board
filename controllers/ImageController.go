package controllers

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type ImageController interface {
    UploadImage(c *gin.Context)
    DeleteImage(c *gin.Context)
}

type ImageControllerImpl struct {}

func NewImageControllerImpl() ImageController {
    return &ImageControllerImpl{}
}

func (i *ImageControllerImpl) UploadImage(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        log.Println(err)
        c.Status(400)
    }

    filename := uuid.NewString()
    if err := c.SaveUploadedFile(file, "./public/images/"+filename+".png"); err != nil {
        log.Println(err)
        c.Status(400)
    }

    domain := os.Getenv("DOMAIN")
    url := "http://"+domain+"/images/"+filename+".png"

    c.JSON(200, gin.H { 
        "url": url,
        "file": filename+".png",
    })
}

func (i *ImageControllerImpl) DeleteImage(c *gin.Context) {
    requestBody := []string{}
    if err := c.ShouldBind(&requestBody); err != nil {
        c.JSON(400, err.Error())
        return
    }
    errs := []string{}
    for _, filename := range requestBody {
        if err := os.Remove("./public/images/"+filename); err!= nil {
            errs = append(errs, filename)
        }
    }
    if len(errs) > 0 {
        c.JSON(400, gin.H {
            "message": "이미지가 삭제되지 않았습니다.",
            "images": errs,
        })
        return
    }

    c.Status(200)
}
