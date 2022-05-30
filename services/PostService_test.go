package services_test

import (
	"okra_board2/config"
	"okra_board2/models"
	"okra_board2/repositories"
	"okra_board2/services"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostService(t *testing.T) {
    db, err := config.InitDBConnection()
    if err != nil { assert.Error(t, err) }
    postRepo := repositories.NewPostRepositoryImpl(db)
    s := services.NewPostServiceImpl(postRepo)

    posts := make([]models.Post, 5)
    for i := 0; i < 5; i++ {
        index := strconv.Itoa(i+1)
        posts[i] = models.Post {
            BoardID: 1,
            Title: "test title " + index,
            Thumbnail: "test thumbnail " + index,
            Content: "test content " + index,
        }
    }

    // insert
    for i := 0; i < len(posts); i++ {
        if err := s.WritePost(&posts[i]); err != nil {
            assert.Error(t, err)
        }
    }

    // reset selected posts
    ids := []int{posts[0].PostID, posts[1].PostID, posts[2].PostID, 1}

    nonexistids, err := s.ResetSelectedPosts(&ids)
    assert.Equal(t, err, gorm.ErrRecordNotFound)
    assert.Equal(t, nonexistids, []int{1})

    ids = []int{posts[0].PostID, posts[1].PostID, posts[2].PostID}

    nonexistids, err = s.ResetSelectedPosts(&ids)
    assert.Equal(t, len(nonexistids), 0)
    assert.Equal(t, err, nil)

    // get selected thumbnails
    thumbnails := s.GetSelectedThumbnails()
    assert.Equal(t, 3, len(thumbnails))

    // delete
    for i := 0; i < len(posts); i++ {
        if err := s.DeletePost(posts[i].PostID); err != nil {
            assert.Error(t, err)
        }
    }

}
