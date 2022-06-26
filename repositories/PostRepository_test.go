package repositories_test

import (
	"errors"
	"okra_board2/config"
	"okra_board2/models"
	"okra_board2/repositories"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostCRUD(t *testing.T) {

    conf, err := config.LoadConfigTest()
    if err != nil { assert.Error(t, err) }

    db, err := config.InitDBConnection(conf)
    if err != nil { assert.Error(t, err) }
    r := repositories.NewPostRepositoryImpl(db)

    sqlDB, err := db.DB()
    defer sqlDB.Close()

    posts := make([]models.Post, 5)
    for i := 0; i < 5; i++ {
        index := strconv.Itoa(i+1)
        posts[i] = models.Post {
            BoardID: 1,
            Title: "test title " + index,
            Thumbnail: "test thumbnail " + index,
            Content: "test content " + index,
            Tags: []models.PostTag {
                { Name: "Tag test 1" },
                { Name: "Tag test 2" },
                { Name: "Tag test 3" },
            },
        }
    }

    // insert
    for i := 0; i < len(posts); i++ {
        if _, err := r.InsertPost(&posts[i]); err != nil {
            assert.Error(t, err)
        }
    }

    // update
    err = r.UpdatePost(&models.Post {
        PostID: posts[0].PostID,
        Title: "updated title",
        Status: true,
    })
    if err != nil { t.Error(err) }


    enabled := true
    // select one
    post, err := r.GetPost(nil, posts[0].PostID)
    if err != nil { t.Error(err) } 

    assert.Equal(t, post.Title, "updated title")

    post, err = r.GetPost(&enabled, posts[1].PostID)
    if err != nil {
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            assert.Error(t, errors.New("GetEnabledPost doesn't work correctly"))
        }
    } else {
        assert.Error(t, errors.New("Something wrong"))
    }

    // check exists
    check := r.CheckPostExists(posts[0].PostID)
    assert.Equal(t, true, check)

    // select many
    keyword := "test title 2"
    boardId := 1
    searchResult, count := r.GetPosts(false, nil, 1, 5, &boardId, &keyword)
    assert.Equal(t, 1, count)
    assert.Equal(t, 1, len(searchResult))

    keyword = "test title"
    searchResult, count = r.GetPosts(false, nil, 1, 5, nil, &keyword)
    assert.Equal(t, 4, count)
    assert.Equal(t, 4, len(searchResult))

    keyword = "updated"
    searchResult, count = r.GetPosts(true, nil, 1, 5, nil, &keyword)
    assert.Equal(t, 1, count)
    assert.Equal(t, 1, len(searchResult))

     // update selected post
//    err = r.ResetSelectedPost(&[]int{posts[0].PostID, posts[1].PostID, posts[2].PostID})
//
//    thumbnails := r.GetSelectedThumbnails()
//    assert.Equal(t, 3, len(thumbnails))
//    assert.Equal(t, "test thumbnail 2", thumbnails[1].Thumbnail)

    // delete all posts
    for i := 0; i < len(posts); i++ {
        if err := r.DeletePost(posts[i].PostID); err != nil {
            t.Error(err)
        }
    }

}
