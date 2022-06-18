package controllers

import (
	"log"
	"okra_board2/config"
    "context"
    "os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageControllerImpl2 struct {
    conf        *config.Config
    client      *s3.Client
}

func NewImageControllerImpl2(
    conf *config.Config,
    client *s3.Client,
) ImageController {
    return &ImageControllerImpl2 {
        conf: conf,
        client: client,
    }
}

func (i *ImageControllerImpl2) UploadImage(c *gin.Context) {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        log.Println(err)
        c.Status(400)
        return
    }

    file, err := fileHeader.Open()
    if err != nil {
        log.Println(err)
        c.Status(400)
        return
    }

    filename := uuid.NewString() + ".png"
    uploader := manager.NewUploader(i.client)
    _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput {
        Bucket: aws.String(i.conf.AWS.Bucket),
        Key:    aws.String("images/"+filename),
        Body:   file,
    })
    if err != nil {
        log.Println(err)
        c.Status(400)
        return
    }

    domain := i.conf.AWS.Domain
    url := "https://"+domain+"/images/"+filename
    c.JSON(200, gin.H {
        "url": url,
        "file": filename,
    })
}

func (i *ImageControllerImpl2) DeleteImage(c *gin.Context) {
    requestBody := []string{}
    if err := c.ShouldBind(&requestBody); err != nil {
        c.JSON(400, err.Error())
        return
    }
    errs := []string{}
    for _, filename := range requestBody {
        if filename == os.Getenv("DEFAULT_THUMBNAIL") {
            continue
        }
        _, err := i.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput {
            Bucket: aws.String(i.conf.AWS.Bucket),
            Key:    aws.String("images/"+filename),
        })
        if err != nil { 
            log.Println(err)
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

