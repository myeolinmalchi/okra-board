package controllers

import (
	"okra_board2/services"
    "okra_board2/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController interface {
    WritePost(c *gin.Context)
    UpdatePost(c *gin.Context)
    DeletePost(c *gin.Context)
    GetPost(enabled bool) gin.HandlerFunc
    GetPosts(enabled bool) gin.HandlerFunc
    ResetSelectedPosts(c *gin.Context)
    GetSelectedThumbnails(c *gin.Context)
}

type PostControllerImpl struct {
    postService services.PostService
}

func NewPostControllerImpl(postService services.PostService) PostController {
    return &PostControllerImpl { postService: postService }
}

func (p *PostControllerImpl) GetPosts(enabled bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var err error
        var (
            size, page int
            boardId *int
            keyword *string
        )
        size, err = strconv.Atoi(c.DefaultQuery("size", "15"))
        if err != nil { c.JSON(400, err.Error()); return }

        page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
        if err != nil { c.JSON(400, err.Error()); return }

        if boardIdStr, boardIdExists := c.GetQuery("boarId"); boardIdExists {
            temp, err := strconv.Atoi(boardIdStr)
            if err != nil { c.JSON(400, err.Error()); return }
            boardId = &temp
        } else {
            boardId = nil
        }
        
        if keywordStr, keywordExists := c.GetQuery("keyword"); keywordExists {
            keyword = &keywordStr
        } else {
            keyword = nil
        }

        posts, count := p.postService.GetPosts(enabled, page, size, boardId, keyword)
        c.IndentedJSON(200, gin.H {
            "nowPage": page,
            "pageCount": int(count / size) + 1,
            "pageSize": size,
            "posts": posts,
        })
    }
}

func (p *PostControllerImpl) GetPost(enabled bool) gin.HandlerFunc {
    return func(c *gin.Context) {

        var err error
        var postId int

        postId, err = strconv.Atoi(c.Param("postId"))
        if err != nil { c.JSON(400, err.Error()); return }
        
        post, err := p.postService.GetPost(enabled, postId)
        if err == gorm.ErrRecordNotFound { c.Status(404); return }

        c.IndentedJSON(200, post)

    }
}

func (p *PostControllerImpl) WritePost(c *gin.Context) {

    requestBody := &models.Post{}

    if err := c.ShouldBind(requestBody); err != nil {
        c.JSON(400, err.Error())
        return
    } 
    postId, result, err := p.postService.WritePost(requestBody)
    if result != nil {
        c.JSON(422, result)
        return
    }
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    c.JSON(200, gin.H {
        "postId": postId,
    })
    
    
}

func (p *PostControllerImpl) UpdatePost(c *gin.Context) {

    postId, err := strconv.Atoi(c.Param("postId"))
    if err != nil { c.JSON(400, err.Error()); return }
    requestBody := &models.Post{}

    if err := c.ShouldBind(requestBody); err != nil {
        c.JSON(400, err.Error())
        return
    } 
    requestBody.PostID = postId
    result, err := p.postService.UpdatePost(requestBody)
    if result != nil {
        c.JSON(422, result)
        return
    }
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    c.Status(200)
}

func (p *PostControllerImpl) DeletePost(c *gin.Context) {
    var err error
    var postId int
    
    postId, err = strconv.Atoi(c.Param("postId"))

    if err != nil { c.JSON(400, err.Error()); return }

    err = p.postService.DeletePost(postId)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.Status(404)
        } else {
            c.JSON(400, err.Error())
        }
    } else {
        c.Status(200)
    }
}

func (p *PostControllerImpl) ResetSelectedPosts(c *gin.Context) {
    requestBody := &[]int{}
    if err := c.ShouldBind(requestBody); err != nil {
        c.JSON(400, err.Error())
        return
    }
    ids, err := p.postService.ResetSelectedPosts(requestBody)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(404, ids)
        } else {
            c.JSON(400, err.Error())
        }
    } else {
        c.Status(200)
    }
}

func (p *PostControllerImpl) GetSelectedThumbnails(c *gin.Context) {
    thumbnails := p.postService.GetSelectedThumbnails()
    c.IndentedJSON(200, thumbnails)
}

