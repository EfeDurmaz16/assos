package handlers

import (
	"assos/api-gateway/internal/models"
	"assos/api-gateway/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VideoHandler struct {
	videoService *services.VideoService
}

func NewVideoHandler(videoService *services.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (h *VideoHandler) GetVideos(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	channelID := c.Query("channel_id")

	if channelID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "channel_id query parameter is required",
		})
	}

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

	videos, err := h.videoService.GetVideosByChannelID(channelID, userID, limit, offset)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to get videos",
		})
	}

	return c.JSON(fiber.Map{
		"videos": videos,
		"count":  len(videos),
		"limit":  limit,
		"offset": offset,
	})
}

func (h *VideoHandler) CreateVideo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req models.CreateVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.ChannelID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "channel_id is required",
		})
	}

	video, err := h.videoService.CreateVideo(userID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to create video",
		})
	}

	return c.Status(201).JSON(video)
}

func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	videoID := c.Params("id")

	video, err := h.videoService.GetVideoByID(videoID, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Video not found",
		})
	}

	return c.JSON(video)
}

func (h *VideoHandler) UpdateVideo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	videoID := c.Params("id")

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	video, err := h.videoService.UpdateVideo(videoID, userID, updates)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to update video",
		})
	}

	return c.JSON(video)
}

func (h *VideoHandler) DeleteVideo(c *fiber.Ctx) error {
	// For now, we'll just update the status to 'deleted' rather than hard delete
	userID := c.Locals("user_id").(string)
	videoID := c.Params("id")

	updates := map[string]interface{}{
		"status": "deleted",
	}

	_, err := h.videoService.UpdateVideo(videoID, userID, updates)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to delete video",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Video deleted successfully",
	})
}

func (h *VideoHandler) ProcessVideo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	videoID := c.Params("id")

	err := h.videoService.ProcessVideo(videoID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to start video processing",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Video processing started",
		"video_id": videoID,
	})
}

func (h *VideoHandler) GetDashboard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// This is a placeholder for dashboard analytics
	// In a real implementation, you would aggregate data from multiple sources
	
	return c.JSON(fiber.Map{
		"user_id": userID,
		"stats": fiber.Map{
			"total_videos": 0,
			"processing_videos": 0,
			"published_videos": 0,
			"total_views": 0,
			"total_revenue": 0,
		},
		"recent_videos": []interface{}{},
		"performance_metrics": fiber.Map{
			"avg_ctr": 0.0,
			"avg_retention": 0.0,
			"avg_rpm": 0.0,
		},
	})
}

func (h *VideoHandler) GetPerformance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	videoID := c.Params("video_id")

	// Verify user owns the video
	video, err := h.videoService.GetVideoByID(videoID, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Video not found",
		})
	}

	// This is a placeholder for performance analytics
	// In a real implementation, you would fetch data from YouTube Analytics API
	
	return c.JSON(fiber.Map{
		"video_id": videoID,
		"title": video.Title,
		"performance": fiber.Map{
			"views": 0,
			"likes": 0,
			"comments": 0,
			"shares": 0,
			"watch_time": 0,
			"ctr": 0.0,
			"retention": 0.0,
			"rpm": 0.0,
		},
		"analytics": fiber.Map{
			"traffic_sources": []interface{}{},
			"demographics": []interface{}{},
			"engagement_over_time": []interface{}{},
		},
	})
}