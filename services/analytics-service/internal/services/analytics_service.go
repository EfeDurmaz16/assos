package services

import (
	"assos/analytics-service/internal/models"
	"database/sql"
	"log"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (s *AnalyticsService) GetDashboardStats(userID string) (*models.DashboardStats, error) {
	log.Printf("Fetching dashboard stats for user %s", userID)

	// In a real implementation, these would be complex queries
	// aggregating data from the 'videos' and 'performance_data' tables.
	// For now, we'll return mock data.
	stats := &models.DashboardStats{
		TotalVideos:      10,
		ProcessingVideos: 2,
		PublishedVideos:  8,
		TotalViews:       123456,
		TotalRevenue:     789.01,
		AvgCTR:           5.6,
		AvgRetention:     45.2,
		AvgRPM:           1.5,
	}

	return stats, nil
}

func (s *AnalyticsService) GetVideoPerformance(videoID string, userID string) (*models.VideoPerformance, error) {
	log.Printf("Fetching performance for video %s for user %s", videoID, userID)

	// First, verify the user owns the video (placeholder query)
	var ownerID string
	err := s.db.QueryRow("SELECT user_id FROM videos v JOIN channels c ON v.channel_id = c.id WHERE v.id = $1", videoID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	// In a real implementation, you would fetch data from YouTube Analytics API
	// or from a dedicated performance data table. For now, we'll return mock data.
	performance := &models.VideoPerformance{
		VideoID:   videoID,
		Title:     "Sample Video Title",
		Views:     12345,
		Likes:     678,
		Comments:  90,
		Shares:    12,
		WatchTime: 54321.0,
		CTR:       8.1,
		Retention: 55.3,
		RPM:       1.8,
	}

	return performance, nil
}

func (s *AnalyticsService) GetRecentVideos(userID string) ([]models.RecentVideo, error) {
	log.Printf("Fetching recent videos for user %s", userID)

	// This would be a query to get the last 5-10 videos
	// For now, we'll return mock data.
	recentVideos := []models.RecentVideo{
		{ID: "uuid-1", Title: "First Recent Video", Status: "published"},
		{ID: "uuid-2", Title: "Second Recent Video", Status: "published"},
		{ID: "uuid-3", Title: "Third Recent Video", Status: "processing"},
	}

	return recentVideos, nil
}