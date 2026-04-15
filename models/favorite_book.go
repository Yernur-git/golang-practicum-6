package models

import "time"

type FavoriteBook struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	BookID    uint      `gorm:"primaryKey" json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}
