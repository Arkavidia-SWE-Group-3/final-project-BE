package handlers

import (
	"Go-Starter-Template/pkg/post"

	"github.com/go-playground/validator/v10"
)

type (
	PostHandler interface {
	}
	postHandler struct {
		PostService post.PostService
		Validator   *validator.Validate
	}
)

func NewPostHandler(postService post.PostService, validator *validator.Validate) PostHandler {
	return &postHandler{
		PostService: postService,
		Validator:   validator,
	}
}
