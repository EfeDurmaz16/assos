package models

import "time"

type DashboardStats struct {
	TotalVideos       int     `json:"total_videos"`
	ProcessingVideos  int     `json:"processing_videos"`
	PublishedVideos   int     `json:"published_videos"`
	TotalViews        int64   `json:"total_views"`
	TotalRevenue      float64 `json:"total_revenue"`
	AvgCTR            float64 `json:"avg_ctr"`
	AvgRetention      float64 `json:"avg_retention"`
	AvgRPM            float64 `json:"avg_rpm"`
}

type VideoPerformance struct {
	VideoID             string    `json:"video_id"`
	Title               string    `json:"title"`
	Views               int64     `json:"views"`
	Likes               int64     `json:"likes"`
	Comments            int64     `json:"comments"`
	Shares              int64     `json:"shares"`
	WatchTime           float64   `json:"watch_time"`
	CTR                 float64   `json:"ctr"`
	Retention           float64   `json:"retention"`
	RPM                 float64   `json:"rpm"`
	LastUpdated         time.Time `json:"last_updated"`
}

type RecentVideo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	PublishedAt time.Time `json:"published_at"`
}