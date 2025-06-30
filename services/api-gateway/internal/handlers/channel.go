package handlers

import (
	"assos/api-gateway/internal/models"
	"assos/api-gateway/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ChannelHandler struct {
	channelService *services.ChannelService
}

func NewChannelHandler(channelService *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
	}
}

func (h *ChannelHandler) GetChannels(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	channels, err := h.channelService.GetChannelsByUserID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get channels",
		})
	}

	return c.JSON(fiber.Map{
		"channels": channels,
		"count":    len(channels),
	})
}

func (h *ChannelHandler) CreateChannel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var channel models.Channel
	if err := c.BodyParser(&channel); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if channel.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Channel name is required",
		})
	}

	createdChannel, err := h.channelService.CreateChannel(userID, &channel)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to create channel",
		})
	}

	return c.Status(201).JSON(createdChannel)
}

func (h *ChannelHandler) GetChannel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	channelID := c.Params("id")

	channel, err := h.channelService.GetChannelByID(channelID, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Channel not found",
		})
	}

	return c.JSON(channel)
}

func (h *ChannelHandler) UpdateChannel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	channelID := c.Params("id")

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	channel, err := h.channelService.UpdateChannel(channelID, userID, updates)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update channel",
		})
	}

	return c.JSON(channel)
}

func (h *ChannelHandler) DeleteChannel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	channelID := c.Params("id")

	err := h.channelService.DeleteChannel(channelID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to delete channel",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Channel deleted successfully",
	})
}