package post

import (
	"Go-Starter-Template/entities"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PostRepository interface {
		CreatePost(ctx context.Context, post entities.Post) error
		UpdatePost(ctx context.Context, post entities.Post) error
		DeletePost(ctx context.Context, postID uuid.UUID) error
		GetPostByID(ctx context.Context, postID uuid.UUID) (entities.Post, error)
		GetFeed(ctx context.Context) ([]entities.Post, error)
	}

	postRepository struct {
		db *gorm.DB
	}
)

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post entities.Post) error {
	if err := r.db.WithContext(ctx).Create(&post).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepository) UpdatePost(ctx context.Context, post entities.Post) error {
	if err := r.db.WithContext(ctx).Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepository) DeletePost(ctx context.Context, postID uuid.UUID) error {

	if err := r.db.WithContext(ctx).Delete(&entities.Post{}, postID).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (entities.Post, error) {
	var post entities.Post
	if err := r.db.WithContext(ctx).First(&post, "id = ?", postID).Error; err != nil {
		return entities.Post{}, err
	}
	return post, nil
}

func (r *postRepository) GetFeed(ctx context.Context) ([]entities.Post, error) {
	var posts []entities.Post
	if err := r.db.WithContext(ctx).Preload("User").Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
