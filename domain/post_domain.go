package domain

import (
	"errors"
	"mime/multipart"
)

var (
	MessageFailedCreatePost = "Failed to create post"
	MessageFailedUpdatePost = "Failed to update post"
	MessageFailedDeletePost = "Failed to delete post"
	MessageFailedGetFeed    = "Failed to get feed"

	MessageSuccessCreatePost = "Successfully create post"
	MessageSuccessUpdatePost = "Successfully update post"
	MessageSuccessDeletePost = "Successfully delete post"
	MessageSuccessGetFeed    = "Successfully get feed"

	ErrCreatePost   = errors.New("failed to create post")
	ErrUpdatePost   = errors.New("failed to update post")
	ErrDeletePost   = errors.New("failed to delete post")
	ErrPostNotFound = errors.New("post not found")
	ErrGetFeed      = errors.New("failed to get feed")
)

type (
	CreatePostRequest struct {
		Content string                `json:"content" form:"content" validate:"required"`
		Asset   *multipart.FileHeader `json:"asset" form:"asset"`
	}

	UpdatePostRequest struct {
		ID      string                `json:"id" form:"id" validate:"required"`
		Content string                `json:"content" form:"content" validate:"required"`
		Asset   *multipart.FileHeader `json:"asset" form:"asset"`
	}

	PostResponse struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Headline       string `json:"headline"`
		ProfilePicture string `json:"profile_picture"`
		Content        string `json:"content"`
		Asset          string `json:"asset"`
	}
)
