package handlers

import (
	"assos/api-gateway/internal/models"
	"assos/api-gateway/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

func (h *AIHandler) GetAgents(c *fiber.Ctx) error {
	agents, err := h.aiService.GetAgents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get agents",
		})
	}

	return c.JSON(fiber.Map{
		"agents": agents,
		"count":  len(agents),
	})
}

func (h *AIHandler) CreateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	agentID := c.Params("id")

	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.VideoID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "video_id is required",
		})
	}

	if req.TaskType == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "task_type is required",
		})
	}

	// Set default priority if not provided
	if req.Priority == 0 {
		req.Priority = 5
	}

	task, err := h.aiService.CreateTask(agentID, userID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	return c.Status(201).JSON(task)
}

func (h *AIHandler) GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// Parse pagination parameters
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	tasks, err := h.aiService.GetTasksByUserID(userID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get tasks",
		})
	}

	return c.JSON(fiber.Map{
		"tasks":  tasks,
		"count":  len(tasks),
		"limit":  limit,
		"offset": offset,
	})
}

func (h *AIHandler) GetTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")

	task, err := h.aiService.GetTaskByID(taskID, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.JSON(task)
}