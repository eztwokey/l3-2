package models

import "time"

type Link struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type Click struct {
	ID        int       `json:"id"`
	LinkID    int       `json:"link_id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	ClickedAt time.Time `json:"clicked_at"`
}

type CreateLinkRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code"`
}

type CreateLinkResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Analytics struct {
	ShortCode   string          `json:"short_code"`
	OriginalURL string          `json:"original_url"`
	TotalClicks int             `json:"total_clicks"`
	ByDay       []DayStat       `json:"by_day"`
	ByMonth     []MonthStat     `json:"by_month"`
	ByUserAgent []UserAgentStat `json:"by_user_agent"`
}

type DayStat struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

type MonthStat struct {
	Month  string `json:"month"`
	Clicks int    `json:"clicks"`
}

type UserAgentStat struct {
	UserAgent string `json:"user_agent"`
	Clicks    int    `json:"clicks"`
}
