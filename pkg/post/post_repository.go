package post

import (
	"gorm.io/gorm"
)

type (
	PostRepository interface {
	}

	postRepository struct {
		db *gorm.DB
	}
)

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}
