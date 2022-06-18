package repositories

import (
	"okra_board2/models"

	"gorm.io/gorm"
)

type PostRepository interface {

    // Select Post and returns with error
    GetPost(postId int)             (post *models.Post, err error)

    // Select Post and if it is unabled, returns error.
    GetEnabledPost(postId int)      (post *models.Post, err error)

    // Insert Post and returns error
    InsertPost(post *models.Post)   (postId int, err error)

    // Update Post and returns error
    UpdatePost(post *models.Post)   (err error)

    // Delete Post and returns error
    DeletePost(postId int)          (err error)

    // Select posts with pagination and optional condition
    // page, size: must contain. parameters for pagination
    // boardId: optional. if nil, select from all boards.
    // keyword: optional. if nil, select all title posts.
    GetPosts(page, size int, boardId *int, keyword *string)         (posts []models.Post, count int)

    // Select enabled posts with pagination and optional condition
    // page, size: must contained. parameters for pagination.
    // boardId: optional. if nil, select from all boards.
    // keyword: optional. if nil, select all title posts.
    GetEnabledPosts(page, size int, boardId *int, keyword *string)  (posts []models.Post, count int)

    // Select posts with pagination, order and optional condition
    // enabled: if true, returns posts which status is true.
    // page, size: must contained. parameters for pagination.
    // boardId: optional. if nil, select from all boards.
    // keyword: optional. if nil, select all title posts.
    // orderBy: order.
    GetPostsOrderBy(
        enabled bool,
        page, size int,
        boardId *int,
        keyword *string,
        orderBy ... string,
    )                               (posts []models.Post, count int)

    // 홈페이지의 메인 화면에 썸네일을 출력 할 게시물들을 재설정한다.
    ResetSelectedPost(ids *[]int)   (err error)

    // Select selected thumbnails.
    GetSelectedThumbnails()         (thumbanils []models.Thumbnail)

    // 게시물이 존재하는지 확인한다.
    CheckPostExists(postId int)     (exists bool)

}

type PostRepositoryImpl struct {
    db *gorm.DB
}

func NewPostRepositoryImpl(db *gorm.DB) PostRepository {
    return &PostRepositoryImpl{ db: db }
}

func (r *PostRepositoryImpl) GetPost(postId int) (post *models.Post, err error) {
    post = &models.Post{}
    err = r.db.Where(&models.Post{PostID:postId}).First(post).Error
    return
}

func (r *PostRepositoryImpl) GetEnabledPost(postId int) (post *models.Post, err error) {
    post = &models.Post{}
    err = r.db.Where(&models.Post{PostID:postId, Status:true}).First(post).Error
    return
}

func (r *PostRepositoryImpl) InsertPost(post *models.Post) (postId int, err error) {
    err = r.db.Create(post).Error
    postId = post.PostID
    return
}

func (r *PostRepositoryImpl) UpdatePost(post *models.Post) (err error) {
    return r.db.UpdateColumns(post).Error
}

func (r *PostRepositoryImpl) DeletePost(postId int) (err error) {
    return r.db.Delete(&models.Post{}, "post_id = ?", postId).Error
}

func (r *PostRepositoryImpl) GetPosts(
    page, size int, 
    boardId *int, 
    keyword *string,
) (posts []models.Post, count int) {
    query := r.db.Model(&models.Post{}).Omit("Content")
    if boardId != nil {
        query = query.Where("board_id = ?", boardId)
    }
    if keyword != nil {
        query = query.Where("title like ?", "%"+*keyword+"%")
    }
    r.db.Table("(?) as a", query).Select("count(*)").Find(&count)
    query.Order("post_id desc").Limit(size).Offset((page-1)*size).Find(&posts)
    return
}

func (r *PostRepositoryImpl) GetEnabledPosts(
    page, size int, 
    boardId *int, 
    keyword *string,
) (posts []models.Post, count int) {
    query := r.db.Model(&models.Post{}).Omit("Content").Where("status = ?", true)
    if boardId != nil {
        query = query.Where("board_id = ?", boardId)
    }
    if keyword != nil {
        query = query.Where("title like ?", "%"+*keyword+"%")
    }
    r.db.Table("(?) as a", query).Select("count(*)").Find(&count)
    query.Order("post_id desc").Limit(size).Offset((page-1)*size).Find(&posts)
    return
}

func (r *PostRepositoryImpl) GetPostsOrderBy(
    enabled bool,
    page, size int,
    boardId *int,
    keyword *string,
    orderBy ... string,
) (posts[]models.Post, count int) {
    query := r.db.Model(&models.Post{}).Omit("Content")
    if enabled { 
        query = query.Where("status = ?", true) 
    }
    if boardId != nil { 
        query = query.Where("board_id = ?", boardId) 
    }
    if keyword != nil { 
        query = query.Where("title like ?", "%"+*keyword+"%") 
    }
    r.db.Table("(?) as a", query).Select("count(*)").Find(&count)
    for _, order := range orderBy {
        query = query.Order(order)
    }
    query.Limit(size).Offset((page-1)*size).Find(&posts)
    return
}

func (r *PostRepositoryImpl) ResetSelectedPost(ids *[]int) (err error) {
    return r.db.Transaction(func(tx *gorm.DB) (err error) {
        err = tx.Model(&models.Post{}).Where("selected = ?", true).Update("selected", false).Error
        if err != nil { return }
        err = tx.Model(&models.Post{}).Where(ids).Update("selected", true).Error
        if err != nil { return }
        return nil
    })
}

func (r *PostRepositoryImpl) GetSelectedThumbnails() (thumbnails []models.Thumbnail) {
    r.db.Table("posts").Where("selected = ? ", true).Order("post_id desc").Find(&thumbnails)
    return
}

func (r *PostRepositoryImpl) CheckPostExists(postId int) (exists bool) {
    r.db.Model(&models.Post{}).
        Select("count(*) > 0").
        Where("post_id = ?", postId).
        Find(&exists)
    return
}
