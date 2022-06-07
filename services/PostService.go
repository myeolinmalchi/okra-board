package services

import (
	"okra_board2/models"
	"okra_board2/repositories"

	"gorm.io/gorm"
)

type PostService interface {

    // 게시물을 작성하고 postId와 유효성 검사 결과 및 에러를 반환한다.
    // post.Thumbnail이 비어있을 경우 "default_thumbnail.png"로 설정한다.
    WritePost(post *models.Post)    (postId int, result *models.PostValidationResult, err error)

    // 게시물을 업데이트하고 유효성 검사 결과와 에러를 반환한다.
    // post.Thumbnail이 비어있을 경우 "default_thumbnail.png"로 설정한다.
    UpdatePost(post *models.Post)   (result *models.PostValidationResult, err error)

    // 게시물을 삭제하고 에러를 반환한다.
    DeletePost(postId int)          (err error)

    // 게시글을 불러온다.
    // enabled 속성이 true일 경우, 
    // status 열이 false인 게시물에 대하여 
    // RecordNotFound 에러를 반환한다.
    GetPost(enabled bool, postId int)             (post *models.Post, err error)
    
    // 조건에 부합하는 게시글의 개수와 함께 게시글 배열을 반환한다.
    // enabled 속성이 true일 경우, status 열이 true인 게시글만을 불러온다.
    // page, size는 페이지네이션을 위한 속성이다.
    // boardId 속성이 nil일 경우 전체 게시판에서 게시글을 검색한다.
    // keyword 속성이 nil이 아닐 경우 제목에 keyword가 포함된 게시글만을 검색한다.
    GetPosts(
        enabled bool,
        page, size int,
        boardId *int,
        keyword *string,
    )                               (posts []models.Post, count int)

    // selected colunm이 true인 게시글들의 썸네일 및 제목 정보를 불러온다.
    GetSelectedThumbnails()         (thumbnaiils []models.Thumbnail)

    // selected column이 true인 게시물을 재설정한다.
    // 전달받은 id 목록 중 존재하지 않는 게시물이 있을 경우
    // 해당 id 리스트를 gorm.ErrRecordNotFound와 함께 반환한다.
    ResetSelectedPosts(ids *[]int)  ([]int, error)

}

type PostServiceImpl struct {
    postRepo repositories.PostRepository
}

func NewPostServiceImpl(
    postRepo repositories.PostRepository,
) PostService {
    return &PostServiceImpl{
        postRepo: postRepo,
    }
}

func (r *PostServiceImpl) checkContent(content string) *string {
    var msg string
    if content == "" {
        msg = "내용을 입력하세요."
    } else {
        return nil
    }
    return &msg
}

func (r *PostServiceImpl) checkTitle(title string) *string {
    var msg string 
    if title == "" {
        msg  = "제목을 입력하세요."
    } else {
        return nil
    }
    return &msg
}

func (r *PostServiceImpl) checkThumbnail(thumbnail string) *string {
    var msg string 
    if thumbnail == "" {
        msg = ""
    } else {
        return nil
    }
    return &msg
}

func (r *PostServiceImpl) postValidation(post *models.Post) *models.PostValidationResult {
    if thumbnailCheck := r.checkThumbnail(post.Thumbnail); thumbnailCheck != nil {
        post.Thumbnail = "default_thumbnail.png"
    }
    result := &models.PostValidationResult {
        Title: r.checkTitle(post.Title),
        Content: r.checkContent(post.Content),
    }
    return result.GetOrNil()
}

func (r *PostServiceImpl) WritePost(post *models.Post) (postId int, result *models.PostValidationResult,  err error) {
    result = r.postValidation(post)
    if result == nil {
        postId, err = r.postRepo.InsertPost(post)
    }
    return
}

func (r *PostServiceImpl) UpdatePost(post *models.Post) (result *models.PostValidationResult, err error) {
    result = r.postValidation(post)
    if result == nil {
        err = r.postRepo.UpdatePost(post)
    }
    return
}

func (r *PostServiceImpl) DeletePost(postId int) (err error) {
    return r. postRepo.DeletePost(postId)
}

func (r *PostServiceImpl) GetPost(enabled bool, postId int) (post *models.Post, err error) {
    if enabled {
        post, err = r.postRepo.GetEnabledPost(postId)
    } else {
        post, err = r.postRepo.GetPost(postId)
    }
    return
}

func (r *PostServiceImpl) GetPosts(
    enabled bool,
    page, size int,
    boardId *int,
    keyword *string,
) (posts []models.Post, count int) {
    if enabled {
        posts, count = r.postRepo.GetEnabledPosts(page, size, boardId, keyword)
    } else {
        posts, count = r.postRepo.GetPosts(page, size, boardId, keyword)
    }
    return
}

func (r *PostServiceImpl) GetSelectedThumbnails() (thumbnails []models.Thumbnail){
    return r.postRepo.GetSelectedThumbnails()
}

func (r *PostServiceImpl) ResetSelectedPosts(ids *[]int) ([]int, error) {
    var nonexistids []int
    for _, id := range *ids {
        if !r.postRepo.CheckPostExists(id) {
            nonexistids = append(nonexistids, id)
        }
    }
    if len(nonexistids) > 0 {
        return nonexistids, gorm.ErrRecordNotFound
    }
    return nil, r.postRepo.ResetSelectedPost(ids)
}
