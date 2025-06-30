package services

import (
	"assos/api-gateway/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ChannelService struct {
	db    *sql.DB
	redis *redis.Client
}

func NewChannelService(db *sql.DB, redis *redis.Client) *ChannelService {
	return &ChannelService{
		db:    db,
		redis: redis,
	}
}

func (s *ChannelService) CreateChannel(userID string, channel *models.Channel) (*models.Channel, error) {
	settingsJSON, _ := json.Marshal(channel.Settings)
	guidelinesJSON, _ := json.Marshal(channel.BrandGuidelines)
	scheduleJSON, _ := json.Marshal(channel.PostingSchedule)

	query := `
		INSERT INTO channels (user_id, youtube_channel_id, name, niche, description, settings, brand_guidelines, posting_schedule)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, youtube_channel_id, name, niche, description, settings, brand_guidelines, posting_schedule, is_active, created_at, updated_at
	`

	var result models.Channel
	err := s.db.QueryRow(query, userID, channel.YoutubeChannelID, channel.Name, channel.Niche, channel.Description, settingsJSON, guidelinesJSON, scheduleJSON).Scan(
		&result.ID,
		&result.UserID,
		&result.YoutubeChannelID,
		&result.Name,
		&result.Niche,
		&result.Description,
		&result.Settings,
		&result.BrandGuidelines,
		&result.PostingSchedule,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return &result, nil
}

func (s *ChannelService) GetChannelsByUserID(userID string) ([]*models.Channel, error) {
	query := `
		SELECT id, user_id, youtube_channel_id, name, niche, description, settings, brand_guidelines, posting_schedule, is_active, created_at, updated_at
		FROM channels
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		var channel models.Channel
		err := rows.Scan(
			&channel.ID,
			&channel.UserID,
			&channel.YoutubeChannelID,
			&channel.Name,
			&channel.Niche,
			&channel.Description,
			&channel.Settings,
			&channel.BrandGuidelines,
			&channel.PostingSchedule,
			&channel.IsActive,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channels = append(channels, &channel)
	}

	return channels, nil
}

func (s *ChannelService) GetChannelByID(id, userID string) (*models.Channel, error) {
	query := `
		SELECT id, user_id, youtube_channel_id, name, niche, description, settings, brand_guidelines, posting_schedule, is_active, created_at, updated_at
		FROM channels
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`

	var channel models.Channel
	err := s.db.QueryRow(query, id, userID).Scan(
		&channel.ID,
		&channel.UserID,
		&channel.YoutubeChannelID,
		&channel.Name,
		&channel.Niche,
		&channel.Description,
		&channel.Settings,
		&channel.BrandGuidelines,
		&channel.PostingSchedule,
		&channel.IsActive,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	return &channel, nil
}

func (s *ChannelService) UpdateChannel(id, userID string, updates map[string]interface{}) (*models.Channel, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	allowedFields := map[string]bool{
		"name":              true,
		"niche":             true,
		"description":       true,
		"youtube_channel_id": true,
		"settings":          true,
		"brand_guidelines":  true,
		"posting_schedule":  true,
	}

	for key, value := range updates {
		if allowedFields[key] {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
			
			// Handle JSONB fields
			if key == "settings" || key == "brand_guidelines" || key == "posting_schedule" {
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
		UPDATE channels 
		SET %s, updated_at = NOW()
		WHERE id = $%d AND user_id = $%d
		RETURNING id, user_id, youtube_channel_id, name, niche, description, settings, brand_guidelines, posting_schedule, is_active, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex, argIndex+1)

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE channels 
			SET %s, updated_at = NOW()
			WHERE id = $%d AND user_id = $%d
			RETURNING id, user_id, youtube_channel_id, name, niche, description, settings, brand_guidelines, posting_schedule, is_active, created_at, updated_at
		`, fmt.Sprintf("%s, %s", setParts[0], setParts[i]), argIndex, argIndex+1)
	}

	args = append(args, id, userID)

	var channel models.Channel
	err := s.db.QueryRow(query, args...).Scan(
		&channel.ID,
		&channel.UserID,
		&channel.YoutubeChannelID,
		&channel.Name,
		&channel.Niche,
		&channel.Description,
		&channel.Settings,
		&channel.BrandGuidelines,
		&channel.PostingSchedule,
		&channel.IsActive,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	return &channel, nil
}

func (s *ChannelService) DeleteChannel(id, userID string) error {
	query := `
		UPDATE channels 
		SET is_active = false, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}