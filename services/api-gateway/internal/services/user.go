package services

import (
	"assos/api-gateway/internal/models"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	db    *sql.DB
	redis *redis.Client
}

func NewUserService(db *sql.DB, redis *redis.Client) *UserService {
	return &UserService{
		db:    db,
		redis: redis,
	}
}

func (s *UserService) CreateUser(email, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, subscription_tier, api_key, is_active, created_at, updated_at
	`

	var user models.User
	err := s.db.QueryRow(query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.SubscriptionTier,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, subscription_tier, api_key, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	var user models.User
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.SubscriptionTier,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, subscription_tier, api_key, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	var user models.User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.SubscriptionTier,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id string, updates map[string]interface{}) (*models.User, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range updates {
		// Only allow specific fields to be updated
		if key == "email" || key == "subscription_tier" {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s, updated_at = NOW()
		WHERE id = $%d
		RETURNING id, email, subscription_tier, api_key, is_active, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE users 
			SET %s, updated_at = NOW()
			WHERE id = $%d
			RETURNING id, email, subscription_tier, api_key, is_active, created_at, updated_at
		`, fmt.Sprintf("%s, %s", setParts[0], setParts[i]), argIndex)
	}

	args = append(args, id)

	var user models.User
	err := s.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.SubscriptionTier,
		&user.APIKey,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}