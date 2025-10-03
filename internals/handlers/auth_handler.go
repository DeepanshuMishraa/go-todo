package handlers

import (
	"strings"

	"github.com/DeepanshuMishraa/gotodo/internals/auth"
	"github.com/DeepanshuMishraa/gotodo/internals/models"
	"github.com/DeepanshuMishraa/gotodo/repository"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req *models.UserRegisterRequest

	if err := c.BodyParser(req); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Invaild Request Body")
	}

	if req.Email == "" || req.Password == "" {
		return RespondWithError(c, fiber.StatusBadRequest, "Email and Password are required")
	}

	if !strings.Contains(req.Email, "@") {
		return RespondWithError(c, fiber.StatusBadRequest, "invalid email format")
	}

	if len(req.Password) < 6 {
		return RespondWithError(c, fiber.StatusBadRequest, "password must be at least 6 characters")
	}

	existingUser, _ := h.userRepo.GetUserByEmail(req.Email)

	if existingUser != nil {
		return RespondWithError(c, fiber.StatusConflict, "user already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "failed to process password")
	}

	user, err := h.userRepo.Create(req.Email, hashedPassword)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "failed to create user")
	}

	token, err := h.jwtManager.GenerateToken(user.Id, user.Email)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "failed to generate token")
	}

	response := fiber.Map{
		"user": models.UserResponse{
			Id:        user.Id,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		"token": token,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}


func (h *AuthHandler) Login(c *fiber.Ctx) error {

	var req models.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "invalid request body")
	}
	
	if req.Email == "" || req.Password == "" {
		return RespondWithError(c, fiber.StatusBadRequest, "email and password are required")
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	if err := auth.CheckPassword(req.Password, user.Password); err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	token, err := h.jwtManager.GenerateToken(user.Id, user.Email)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "failed to generate token")
	}

	response := fiber.Map{
		"user": models.UserResponse{
			Id:        user.Id,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		"token": token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}