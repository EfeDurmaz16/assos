package handlers

import (
	"assos/analytics-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	// In a real app, you would get the user ID from a JWT or other auth mechanism
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	stats, err := h.service.GetDashboardStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get dashboard stats"})
	}

	recentVideos, err := h.service.GetRecentVideos(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get recent videos"})
	}

	return c.JSON(fiber.Map{
		"stats":         stats,
		"recent_videos": recentVideos,
	})
}

func (h *AnalyticsHandler) GetPerformance(c *fiber.Ctx) error {
	videoID := c.Params("video_id")
	userID := c.Query("user_id") // As above, this would come from auth
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	performance, err := h.service.GetVideoPerformance(videoID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get video performance"})
	}
	if performance == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "video not found or access denied"})
	}

	return c.JSON(performance)
}