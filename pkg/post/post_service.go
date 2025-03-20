package post

import (
	jwtService "Go-Starter-Template/pkg/jwt"
)

type (
	PostService interface {
	}

	postService struct {
		postRepository PostRepository
		jwtService     jwtService.JWTService
	}
)

func NewPostService(postRepository PostRepository, jwtService jwtService.JWTService) PostService {
	return &postService{postRepository: postRepository, jwtService: jwtService}
}
