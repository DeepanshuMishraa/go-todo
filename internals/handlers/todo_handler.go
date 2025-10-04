package handlers

import (
	"strconv"

	"github.com/DeepanshuMishraa/gotodo/internals/models"
	"github.com/DeepanshuMishraa/gotodo/repository"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	todoRepo *repository.TodoRepository
}

func NewTodoHandler(todoRepo *repository.TodoRepository) *TodoHandler {
	return &TodoHandler{
		todoRepo: todoRepo,
	}
}

func (h *TodoHandler) Create(c *fiber.Ctx) error {
	userId := c.Locals("userID").(int)

	var req *models.CreateTodoRequest

	if err := c.BodyParser(&req); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Title == "" || req.Description == "" {
		return RespondWithError(c, fiber.StatusBadRequest, "Title and Description are required")
	}

	todo, err := h.todoRepo.CreateTodo(userId, req.Title, req.Description)

	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": todo,
	})
}

func (h *TodoHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	todoID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Invalid todo ID")
	}

	todo, err := h.todoRepo.GetByID(todoID)

	if err != nil {
		return RespondWithError(c, fiber.StatusNotFound, "Todo not found")
	}

	if todo.UserId != userID {
		return RespondWithError(c, fiber.StatusForbidden, "access denied")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": todo,
	})
}

func (h *TodoHandler) Update(c *fiber.Ctx) error {
	userId := c.Locals("userID").(int)

	todoID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Invalid todo ID")
	}

	existingTodo, err := h.todoRepo.GetByID(todoID)

	if err != nil {
		return RespondWithError(c, fiber.StatusNotFound, "Todo not found")
	}

	if existingTodo.UserId != userId {
		return RespondWithError(c, fiber.StatusForbidden, "access denied")
	}

	var req models.UpdateTodoRequest

	if err := c.BodyParser(&req); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	title := existingTodo.Title
	if req.Title != nil {
		title = *req.Title
	}

	description := existingTodo.Description
	if req.Description != nil {
		description = *req.Description
	}

	completed := existingTodo.IsCompleted
	if req.IsCompleted != nil {
		completed = *req.IsCompleted
	}

	updatedTodo, err := h.todoRepo.Update(todoID, &title, &description, &completed)

	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "Failed to update todo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": updatedTodo,
	})
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	todoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "invalid todo ID")
	}

	todo, err := h.todoRepo.GetByID(todoID)
	if err != nil {
		return RespondWithError(c, fiber.StatusNotFound, "todo not found")
	}

	if todo.UserId != userID {
		return RespondWithError(c, fiber.StatusForbidden, "access denied")
	}

	if err := h.todoRepo.Delete(todoID); err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "failed to delete todo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "todo deleted successfully",
	})
}
