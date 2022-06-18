package models

import "time"

type Post struct {
    PostID      int         `json:"postId,omitempty" gorm:"primaryKey;<-:false"`
    BoardID     int         `json:"boardId"`
    Title       string      `json:"title"`
    Thumbnail   string      `json:"thumbnail"`
    Content     string      `json:"content,omitempty"`
    AddedDate   time.Time   `json:"addedDate,omitempty" gorm:"->"`
    Status      bool        `json:"status"`
    Selected    bool        `json:"selected"`
    Views       int         `json:"views"`
}

// Response Only
type PostValidationResult struct {
    Title       *string     `json:"title,omitempty"`
    Thumbnail   *string     `json:"thumbnail,omitempty"`
    Content     *string     `json:"content,omitempty"`
}

func (result *PostValidationResult) GetOrNil() *PostValidationResult {
    if result.Title == nil && result.Thumbnail == nil && result.Content == nil {
        return nil
    }
    return result
}

// Response Only
type Thumbnail struct {
    PostID      int         `json:"postId"`
    Title       string      `json:"title"`
    Thumbnail   string      `json:"thumbnail"`
}
