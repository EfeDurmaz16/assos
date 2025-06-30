package services

import (
	"assos/api-gateway/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type AIService struct {
	db   *sql.DB
	redis *redis.Client
	nats *nats.Conn
}

func NewAIService(db *sql.DB, redis *redis.Client, nats *nats.Conn) *AIService {
	return &AIService{
		db:   db,
		redis: redis,
		nats: nats,
	}
}

func (s *AIService) GetAgents() ([]*models.AIAgent, error) {
	query := `
		SELECT id, name, type, configuration, performance_metrics, is_active, created_at, updated_at
		FROM ai_agents
		WHERE is_active = true
		ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}
	defer rows.Close()

	var agents []*models.AIAgent
	for rows.Next() {
		var agent models.AIAgent
		err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Type,
			&agent.Configuration,
			&agent.PerformanceMetrics,
			&agent.IsActive,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		agents = append(agents, &agent)
	}

	return agents, nil
}

func (s *AIService) CreateTask(agentID, userID string, req *models.CreateTaskRequest) (*models.AgentTask, error) {
	// Verify user owns the video
	var videoExists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM videos v 
			JOIN channels c ON v.channel_id = c.id 
			WHERE v.id = $1 AND c.user_id = $2
		)
	`, req.VideoID, userID).Scan(&videoExists)
	
	if err != nil || !videoExists {
		return nil, fmt.Errorf("video not found or not owned by user")
	}

	// Verify agent exists
	var agentExists bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ai_agents WHERE id = $1 AND is_active = true)", agentID).Scan(&agentExists)
	if err != nil || !agentExists {
		return nil, fmt.Errorf("agent not found")
	}

	inputDataJSON, _ := json.Marshal(req.InputData)

	query := `
		INSERT INTO agent_tasks (agent_id, video_id, task_type, priority, input_data, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		RETURNING id, agent_id, video_id, task_type, priority, status, input_data, output_data, execution_time, created_at, updated_at
	`

	var task models.AgentTask
	err = s.db.QueryRow(query, agentID, req.VideoID, req.TaskType, req.Priority, inputDataJSON).Scan(
		&task.ID,
		&task.AgentID,
		&task.VideoID,
		&task.TaskType,
		&task.Priority,
		&task.Status,
		&task.InputData,
		&task.OutputData,
		&task.ExecutionTime,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Send task to AI service via NATS
	message := map[string]interface{}{
		"task_id":   task.ID,
		"agent_id":  agentID,
		"video_id":  req.VideoID,
		"task_type": req.TaskType,
		"priority":  req.Priority,
		"input_data": req.InputData,
	}

	messageBytes, _ := json.Marshal(message)
	err = s.nats.Publish("ai.task", messageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to queue task: %w", err)
	}

	return &task, nil
}

func (s *AIService) GetTasksByUserID(userID string, limit, offset int) ([]*models.AgentTask, error) {
	query := `
		SELECT t.id, t.agent_id, t.video_id, t.task_type, t.priority, t.status, t.input_data, t.output_data, t.execution_time, t.created_at, t.updated_at
		FROM agent_tasks t
		JOIN videos v ON t.video_id = v.id
		JOIN channels c ON v.channel_id = c.id
		WHERE c.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.AgentTask
	for rows.Next() {
		var task models.AgentTask
		err := rows.Scan(
			&task.ID,
			&task.AgentID,
			&task.VideoID,
			&task.TaskType,
			&task.Priority,
			&task.Status,
			&task.InputData,
			&task.OutputData,
			&task.ExecutionTime,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (s *AIService) GetTaskByID(taskID, userID string) (*models.AgentTask, error) {
	query := `
		SELECT t.id, t.agent_id, t.video_id, t.task_type, t.priority, t.status, t.input_data, t.output_data, t.execution_time, t.created_at, t.updated_at
		FROM agent_tasks t
		JOIN videos v ON t.video_id = v.id
		JOIN channels c ON v.channel_id = c.id
		WHERE t.id = $1 AND c.user_id = $2
	`

	var task models.AgentTask
	err := s.db.QueryRow(query, taskID, userID).Scan(
		&task.ID,
		&task.AgentID,
		&task.VideoID,
		&task.TaskType,
		&task.Priority,
		&task.Status,
		&task.InputData,
		&task.OutputData,
		&task.ExecutionTime,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	return &task, nil
}