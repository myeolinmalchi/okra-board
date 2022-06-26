package repositories

import (
	"okra_board2/models"

	"gorm.io/gorm"
)

type PostRepository interface {

    // posts 테이블에서 게시물 정보를 불러온다.
    // status == nil => 게시물의 status를 구분하지 않고 검색한다.
    // status != nil => 지정된 status의 게시물중에서 검색한다. 
    // 조건에 부합하는 게시글을 찾지 못할 경우 err 반환.
    GetPost(
        status *bool,
        postId int,
    )                               (post *models.Post, err error)

    // 이전 게시물의 post_id와 title 정보를 검색한다.
    // status == nil => 게시물의 status를 구분하지 않고 검색한다.
    // status != nil => 지정된 status의 게시물중에서 검색한다. 
    // 조건에 부합하는 게시글을 찾지 못할 경우 err 반환.
    GetPrevPostInfo(
        status *bool,
        postId int,
    )                               (prevPost *models.PostE, err error)

    // 다음 게시물의 post_id와 title 정보를 검색한다.
    // status == nil => 게시물의 status를 구분하지 않고 검색한다.
    // status != nil => 지정된 status의 게시물중에서 검색한다. 
    // 조건에 부합하는 게시글을 찾지 못할 경우 err 반환.
    GetNextPostInfo(
        status *bool,
        postId int,
    )                               (nextPost *models.PostE, err error)

    // Insert Post and returns error
    InsertPost(post *models.Post)   (postId int, err error)

    // Update Post and returns error
    UpdatePost(post *models.Post)   (err error)

    // Delete Post and returns error
    DeletePost(postId int)          (err error)

    // Select posts with pagination, order and optional condition
    // enabled: if true, returns posts which status is true.
    // page, size: must be contained. parameters for pagination.
    // boardId: optional. if nil, select from all boards.
    // keyword: optional. if nil, select all title posts.
    // orderBy: order.
    GetPosts(
        enabled bool,
        selected *bool,
        page, size int,
        boardId *int,
        titleKeyword *string,
        tagKeyword *string,
        orderBy ... string,
    )                               (posts []models.Post, count int)

    // posts 테이블의 모든 게시글 정보를 불러온다.
    GetAllPosts()                   (posts []models.PostE)

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

func (r *PostRepositoryImpl) GetPost(status *bool, postId int) (post *models.Post, err error) {
    post = &models.Post{}
    query := r.db.Model(&models.Post{}).Preload("Tags", func(db *gorm.DB) *gorm.DB {
        return db.Order("post_tags.name ASC")
    })
    if status != nil {
        query = query.Where("status = ?", *status)
    }
    err = query.Where("post_id = ?", postId).First(post).Error
    return
}

func (r *PostRepositoryImpl) GetPrevPostInfo(
    status *bool,
    postId int,
) (prevPost *models.PostE, err error) {
    prevPost = &models.PostE{}
    query := r.db.Model(&models.Post{})    
    if status != nil {
        query = query.Where("status = ?", *status)
    }
    err = query.Where("post_id < ?", postId).Order("post_id desc").First(prevPost).Error
    return
}

func (r *PostRepositoryImpl) GetNextPostInfo(
    status *bool,
    postId int,
) (nextPost *models.PostE, err error) {
    nextPost = &models.PostE{}
    query := r.db.Model(&models.Post{}).Select("post_id, title")
    if status != nil {
        query = query.Where("status = ?", *status)
    }
    err = query.Where("post_id > ?", postId).Order("post_id asc").First(nextPost).Error
    return
}

func (r *PostRepositoryImpl) InsertPost(post *models.Post) (postId int, err error) {
    err = r.db.Create(post).Error
    postId = post.PostID
    return
}

func (r *PostRepositoryImpl) UpdatePost(post *models.Post) (err error) {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := r.db.Delete(&models.PostTag{}, "post_id = ?", post.PostID).Error; err != nil {
            return err
        }
        if len(post.Tags) > 0  {
            if err := r.db.Create(post.Tags).Error; err != nil {
                return err
            }
        }
        if err := r.db.UpdateColumns(post).Error; err != nil {
            return err
        }
        return nil
    })
}

func (r *PostRepositoryImpl) DeletePost(postId int) (err error) {
    return r.db.Delete(&models.Post{}, "post_id = ?", postId).Error
}

func (r *PostRepositoryImpl) GetPosts(
    enabled bool,
    selected *bool,
    page, size int,
    boardId *int,
    titleKeyword *string,
    tagKeyword *string,
    orderBy ... string,
) (posts[]models.Post, count int) {
    query := r.db.Model(&models.Post{}).Preload("Tags", func(db *gorm.DB) *gorm.DB {
        return db.Order("post_tags.name ASC")
    }).Omit("Content")
    if enabled { 
        query = query.Where("status = ?", true) 
    }
    if selected != nil {
        query = query.Where("selected = ?", selected)
    }
    if boardId != nil { 
        query = query.Where("board_id = ?", boardId) 
    }
    if titleKeyword != nil { 
        query = query.Where("title like ?", "%"+*titleKeyword+"%") 
    }
    if tagKeyword != nil {
        query = query.Joins("INNER JOIN post_tags on post_tags.post_id = posts.post_id").
            Group("posts.post_id").
            Where("post_tags.name like ?", "%"+*tagKeyword+"%")
    }
    r.db.Table("(?) as a", query).Select("count(*)").Find(&count)
    for _, order := range orderBy {
        query = query.Order(order)
    }
    query.Limit(size).Offset((page-1)*size).Find(&posts)
    return
}

func (r *PostRepositoryImpl) GetAllPosts() (posts []models.PostE){
    r.db.Model(&models.Post{}).Find(&posts)
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
