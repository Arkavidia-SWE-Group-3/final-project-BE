package handlers

import (
	"Go-Starter-Template/pkg/post"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"

	"Go-Starter-Template/domain"

	"Go-Starter-Template/internal/api/presenters"
)

type (
	PostHandler interface {
		CreatePost(c *fiber.Ctx) error
		UpdatePost(c *fiber.Ctx) error
		DeletePost(c *fiber.Ctx) error
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

func (h *postHandler) CreatePost(c *fiber.Ctx) error {
	var req domain.CreatePostRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedCreatePost, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedCreatePost, err)
	}

	userID := c.Locals("user_id").(string)
	req.Asset, _ = c.FormFile("asset")

	err := h.PostService.CreatePost(c.Context(), req, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedCreatePost, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusOK, domain.MessageSuccessCreatePost)
}

func (h *postHandler) UpdatePost(c *fiber.Ctx) error {
	var req domain.UpdatePostRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedUpdatePost, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedUpdatePost, err)
	}

	userID := c.Locals("user_id").(string)
	req.Asset, _ = c.FormFile("asset")

	err := h.PostService.UpdatePost(c.Context(), req, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedUpdatePost, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusOK, domain.MessageSuccessUpdatePost)
}

func (h *postHandler) DeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	userID := c.Locals("user_id").(string)

	err := h.PostService.DeletePost(c.Context(), postID, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedDeletePost, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusOK, domain.MessageSuccessDeletePost)
}
