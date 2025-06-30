package services

import (
	"assos/api-gateway/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type VideoService struct {
	db   *sql.DB
	redis *redis.Client
	nats *nats.Conn
}

func NewVideoService(db *sql.DB, redis *redis.Client, nats *nats.Conn) *VideoService {
	return &VideoService{
		db:   db,
		redis: redis,
		nats: nats,
	}
}

func (s *VideoService) CreateVideo(userID string, req *models.CreateVideoRequest) (*models.Video, error) {
	// Verify user owns the channel
	var channelExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1 AND user_id = $2)", req.ChannelID, userID).Scan(&channelExists)
	if err != nil || !channelExists {
		return nil, fmt.Errorf("channel not found or not owned by user")
	}

	metadataJSON, _ := json.Marshal(req.Metadata)

	query := `
		INSERT INTO videos (channel_id, title, description, metadata, status)
		VALUES ($1, $2, $3, $4, 'research')
		RETURNING id, channel_id, title, description, status, metadata, created_at, updated_at
	`

	var video models.Video
	err = s.db.QueryRow(query, req.ChannelID, req.Title, req.Description, metadataJSON).Scan(
		&video.ID,
		&video.ChannelID,
		&video.Title,
		&video.Description,
		&video.Status,
		&video.Metadata,
		&video.CreatedAt,
		&video.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create video: %w", err)
	}

	return &video, nil
}

func (s *VideoService) GetVideosByChannelID(channelID, userID string, limit, offset int) ([]*models.Video, error) {
	// Verify user owns the channel
	var channelExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1 AND user_id = $2)", channelID, userID).Scan(&channelExists)
	if err != nil || !channelExists {
		return nil, fmt.Errorf("channel not found or not owned by user")
	}

	query := `
		SELECT id, channel_id, title, description, status, youtube_video_id, thumbnail_url, video_url, 
			   script, metadata, performance_data, ai_analysis, processing_started_at, processing_completed_at, 
			   published_at, created_at, updated_at
		FROM videos
		WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, channelID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}
	defer rows.Close()

	var videos []*models.Video
	for rows.Next() {
		var video models.Video
		err := rows.Scan(
			&video.ID,
			&video.ChannelID,
			&video.Title,
			&video.Description,
			&video.Status,
			&video.YoutubeVideoID,
			&video.ThumbnailURL,
			&video.VideoURL,
			&video.Script,
			&video.Metadata,
			&video.PerformanceData,
			&video.AIAnalysis,
			&video.ProcessingStartedAt,
			&video.ProcessingCompletedAt,
			&video.PublishedAt,
			&video.CreatedAt,
			&video.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func (s *VideoService) GetVideoByID(id, userID string) (*models.Video, error) {
	query := `
		SELECT v.id, v.channel_id, v.title, v.description, v.status, v.youtube_video_id, v.thumbnail_url, v.video_url, 
			   v.script, v.metadata, v.performance_data, v.ai_analysis, v.processing_started_at, v.processing_completed_at, 
			   v.published_at, v.created_at, v.updated_at
		FROM videos v
		JOIN channels c ON v.channel_id = c.id
		WHERE v.id = $1 AND c.user_id = $2
	`

	var video models.Video
	err := s.db.QueryRow(query, id, userID).Scan(
		&video.ID,
		&video.ChannelID,
		&video.Title,
		&video.Description,
		&video.Status,
		&video.YoutubeVideoID,
		&video.ThumbnailURL,
		&video.VideoURL,
		&video.Script,
		&video.Metadata,
		&video.PerformanceData,
		&video.AIAnalysis,
		&video.ProcessingStartedAt,
		&video.ProcessingCompletedAt,
		&video.PublishedAt,
		&video.CreatedAt,
		&video.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("video not found: %w", err)
	}

	return &video, nil
}

func (s *VideoService) UpdateVideo(id, userID string, updates map[string]interface{}) (*models.Video, error) {
	// Verify user owns the video
	var videoExists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM videos v 
			JOIN channels c ON v.channel_id = c.id 
			WHERE v.id = $1 AND c.user_id = $2
		)
	`, id, userID).Scan(&videoExists)
	
	if err != nil || !videoExists {
		return nil, fmt.Errorf("video not found or not owned by user")
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	allowedFields := map[string]bool{
		"title":        true,
		"description":  true,
		"status":       true,
		"metadata":     true,
		"script":       true,
	}

	for key, value := range updates {
		if allowedFields[key] {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
			
			// Handle JSONB fields
			if key == "metadata" || key == "script" {
				jsonValue, _ := json.Marshal(value)
				args = append(args, string(jsonValue))
			} else {
				args = append(args, value)
			}
			argIndex++
		}
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE videos 
		SET %s, updated_at = NOW()
		WHERE id = $%d
		RETURNING id, channel_id, title, description, status, youtube_video_id, thumbnail_url, video_url, 
				  script, metadata, performance_data, ai_analysis, processing_started_at, processing_completed_at, 
				  published_at, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE videos 
			SET %s, updated_at = NOW()
			WHERE id = $%d
			RETURNING id, channel_id, title, description, status, youtube_video_id, thumbnail_url, video_url, 
					  script, metadata, performance_data, ai_analysis, processing_started_at, processing_completed_at, 
					  published_at, created_at, updated_at
		`, fmt.Sprintf("%s, %s", setParts[0], setParts[i]), argIndex)
	}

	args = append(args, id)

	var video models.Video
	err = s.db.QueryRow(query, args...).Scan(
		&video.ID,
		&video.ChannelID,
		&video.Title,
		&video.Description,
		&video.Status,
		&video.YoutubeVideoID,
		&video.ThumbnailURL,
		&video.VideoURL,
		&video.Script,
		&video.Metadata,
		&video.PerformanceData,
		&video.AIAnalysis,
		&video.ProcessingStartedAt,
		&video.ProcessingCompletedAt,
		&video.PublishedAt,
		&video.CreatedAt,
		&video.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update video: %w", err)
	}

	return &video, nil
}

func (s *VideoService) ProcessVideo(videoID, userID string) error {
	// Verify user owns the video
	var videoExists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM videos v 
			JOIN channels c ON v.channel_id = c.id 
			WHERE v.id = $1 AND c.user_id = $2
		)
	`, videoID, userID).Scan(&videoExists)
	
	if err != nil || !videoExists {
		return fmt.Errorf("video not found or not owned by user")
	}

	// Update video status to processing
	_, err = s.db.Exec(`
		UPDATE videos 
		SET status = 'scripting', processing_started_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, videoID)
	
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}

	// Send message to AI service to start processing
	message := map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
		"action":   "start_processing",
	}

	messageBytes, _ := json.Marshal(message)
	
	err = s.nats.Publish("video.process", messageBytes)
	if err != nil {
		return fmt.Errorf("failed to queue video processing: %w", err)
	}

	return nil
}