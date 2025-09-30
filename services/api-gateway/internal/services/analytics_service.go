package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AnalyticsService struct {
	baseURL string
}

func NewAnalyticsService(baseURL string) *AnalyticsService {
	return &AnalyticsService{baseURL: baseURL}
}

func (s *AnalyticsService) GetDashboard(userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/analytics/dashboard?user_id=%s", s.baseURL, userID)
	return s.makeRequest(url)
}

func (s *AnalyticsService) GetPerformance(videoID, userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/analytics/performance/%s?user_id=%s", s.baseURL, videoID, userID)
	return s.makeRequest(url)
}

func (s *AnalyticsService) makeRequest(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to analytics service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("analytics service returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response from analytics service: %w", err)
	}

	return data, nil
}