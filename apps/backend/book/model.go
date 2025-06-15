package book

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null" validate:"required"`
	Author    string         `json:"author" gorm:"not null" validate:"required"`
	Year      int            `json:"year" gorm:"not null" validate:"required"`
	Genre     string         `json:"genre"`
	ISBN      string         `json:"isbn" gorm:"uniqueIndex"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
