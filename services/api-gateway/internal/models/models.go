package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	}

	return nil
}

// User represents a user in the system
type User struct {
	ID               string    `json:"id" db:"id"`
	Email            string    `json:"email" db:"email"`
	PasswordHash     string    `json:"-" db:"password_hash"`
	SubscriptionTier string    `json:"subscription_tier" db:"subscription_tier"`
	APIKey           string    `json:"api_key" db:"api_key"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Channel represents a YouTube channel
type Channel struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	YoutubeChannelID  string    `json:"youtube_channel_id" db:"youtube_channel_id"`
	Name              string    `json:"name" db:"name"`
	Niche             string    `json:"niche" db:"niche"`
	Description       string    `json:"description" db:"description"`
	Settings          JSONB     `json:"settings" db:"settings"`
	BrandGuidelines   JSONB     `json:"brand_guidelines" db:"brand_guidelines"`
	PostingSchedule   JSONB     `json:"posting_schedule" db:"posting_schedule"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Video represents a video in the system
type Video struct {
	ID                     string     `json:"id" db:"id"`
	ChannelID              string     `json:"channel_id" db:"channel_id"`
	Title                  string     `json:"title" db:"title"`
	Description            string     `json:"description" db:"description"`
	Status                 string     `json:"status" db:"status"`
	YoutubeVideoID         string     `json:"youtube_video_id" db:"youtube_video_id"`
	ThumbnailURL           string     `json:"thumbnail_url" db:"thumbnail_url"`
	VideoURL               string     `json:"video_url" db:"video_url"`
	Script                 JSONB      `json:"script" db:"script"`
	Metadata               JSONB      `json:"metadata" db:"metadata"`
	PerformanceData        JSONB      `json:"performance_data" db:"performance_data"`
	AIAnalysis             JSONB      `json:"ai_analysis" db:"ai_analysis"`
	ProcessingStartedAt    *time.Time `json:"processing_started_at" db:"processing_started_at"`
	ProcessingCompletedAt  *time.Time `json:"processing_completed_at" db:"processing_completed_at"`
	PublishedAt            *time.Time `json:"published_at" db:"published_at"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}

// ContentPipeline represents a stage in the content creation pipeline
type ContentPipeline struct {
	ID             string     `json:"id" db:"id"`
	VideoID        string     `json:"video_id" db:"video_id"`
	Stage          string     `json:"stage" db:"stage"`
	Status         string     `json:"status" db:"status"`
	InputData      JSONB      `json:"input_data" db:"input_data"`
	OutputData     JSONB      `json:"output_data" db:"output_data"`
	ErrorMessage   string     `json:"error_message" db:"error_message"`
	ProcessingTime int        `json:"processing_time" db:"processing_time"`
	AgentUsed      string     `json:"agent_used" db:"agent_used"`
	StartedAt      *time.Time `json:"started_at" db:"started_at"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// AIAgent represents an AI agent in the system
type AIAgent struct {
	ID                 string    `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	Type               string    `json:"type" db:"type"`
	Configuration      JSONB     `json:"configuration" db:"configuration"`
	PerformanceMetrics JSONB     `json:"performance_metrics" db:"performance_metrics"`
	IsActive           bool      `json:"is_active" db:"is_active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// AgentTask represents a task assigned to an AI agent
type AgentTask struct {
	ID            string     `json:"id" db:"id"`
	AgentID       string     `json:"agent_id" db:"agent_id"`
	VideoID       string     `json:"video_id" db:"video_id"`
	TaskType      string     `json:"task_type" db:"task_type"`
	Priority      int        `json:"priority" db:"priority"`
	Status        string     `json:"status" db:"status"`
	InputData     JSONB      `json:"input_data" db:"input_data"`
	OutputData    JSONB      `json:"output_data" db:"output_data"`
	ExecutionTime int        `json:"execution_time" db:"execution_time"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response models
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateVideoRequest struct {
	ChannelID   string `json:"channel_id" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Metadata    JSONB  `json:"metadata"`
}

type CreateTaskRequest struct {
	VideoID   string `json:"video_id" validate:"required"`
	TaskType  string `json:"task_type" validate:"required"`
	Priority  int    `json:"priority"`
	InputData JSONB  `json:"input_data"`
}