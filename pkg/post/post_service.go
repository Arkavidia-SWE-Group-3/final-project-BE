package post

import (
	jwtService "Go-Starter-Template/pkg/jwt"

	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	"context"

	"github.com/google/uuid"
)

type (
	PostService interface {
		CreatePost(ctx context.Context, req domain.CreatePostRequest, userID string) error
		UpdatePost(ctx context.Context, req domain.UpdatePostRequest, userID string) error
		DeletePost(ctx context.Context, postID string, userID string) error
	}

	postService struct {
		postRepository PostRepository
		awsS3          storage.AwsS3
		jwtService     jwtService.JWTService
	}
)

func NewPostService(postRepository PostRepository, awsS3 storage.AwsS3, jwtService jwtService.JWTService) PostService {
	return &postService{postRepository: postRepository, awsS3: awsS3, jwtService: jwtService}
}

func (s *postService) CreatePost(ctx context.Context, req domain.CreatePostRequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	post := entities.Post{
		UserID:  parsedUserID,
		Content: req.Content,
	}

	allowedMimetype := []string{"image/jpeg", "image/jpg", "image/png"}

	if req.Asset != nil {
		objectKey, err := s.awsS3.UploadFile(utils.GenerateRandomFileName(req.Content), req.Asset, "posts", allowedMimetype...)
		if err != nil {
			return domain.ErrUploadFile
		}

		post.Asset = s.awsS3.GetPublicLinkKey(objectKey)
	}

	err = s.postRepository.CreatePost(ctx, post)

	if err != nil {
		return domain.ErrCreatePost
	}

	return nil
}

func (s *postService) UpdatePost(ctx context.Context, req domain.UpdatePostRequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	parsedPostID, err := uuid.Parse(req.ID)

	if err != nil {
		return domain.ErrParseUUID
	}

	post, err := s.postRepository.GetPostByID(ctx, parsedPostID)

	if err != nil {
		return domain.ErrPostNotFound
	}

	if post.UserID != parsedUserID {
		return domain.ErrUserNotAllowed
	}

	post = entities.Post{
		ID:      parsedPostID,
		UserID:  parsedUserID,
		Content: req.Content,
	}

	if req.Asset != nil {
		allowedMimetype := []string{"image/jpeg", "image/jpg", "image/png"}

		objectKey, err := s.awsS3.UploadFile(utils.GenerateRandomFileName(req.Content), req.Asset, "posts", allowedMimetype...)
		if err != nil {
			return domain.ErrUploadFile
		}

		post.Asset = s.awsS3.GetPublicLinkKey(objectKey)
	}

	err = s.postRepository.UpdatePost(ctx, post)

	if err != nil {
		return domain.ErrUpdatePost
	}

	return nil
}

func (s *postService) DeletePost(ctx context.Context, postID string, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	parsedPostID, err := uuid.Parse(postID)

	if err != nil {
		return domain.ErrParseUUID
	}

	post, err := s.postRepository.GetPostByID(ctx, parsedPostID)

	if err != nil {
		return domain.ErrPostNotFound
	}

	if post.UserID != parsedUserID {
		return domain.ErrUserNotAllowed
	}

	err = s.postRepository.DeletePost(ctx, parsedPostID)

	if err != nil {
		return domain.ErrDeletePost
	}

	return nil
}
