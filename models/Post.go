package models

import "time"

type Post struct {
    PostID      int         `json:"postId,omitempty" gorm:"primaryKey;<-:false"`
    BoardID     int         `json:"boardId"`
    Title       string      `json:"title"`
    Thumbnail   string      `json:"thumbnail"`
    Content     string      `json:"content"`
    AddedDate   time.Time   `json:"addedDate,omitempty" gorm:"->"`
    Status      bool        `json:"status"`
    Selected    bool        `json:"selected"`
    Views       int         `json:"views"`
}

// Response Only
type Thumbnail struct {
    PostID      int         `json:"postId"`
    Title       string      `json:"title"`
    Thumbnail   string      `json:"thumbnail"`
}
