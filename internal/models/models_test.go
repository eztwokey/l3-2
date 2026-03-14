package models

import (
	"encoding/json"
	"testing"
)

func TestCreateLinkRequest_JSON(t *testing.T) {
	input := `{"url":"https://example.com","custom_code":"test"}`

	var req CreateLinkRequest
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if req.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got %q", req.URL)
	}
	if req.CustomCode != "test" {
		t.Errorf("expected CustomCode 'test', got %q", req.CustomCode)
	}
}

func TestCreateLinkResponse_JSON(t *testing.T) {
	resp := CreateLinkResponse{
		ShortCode:   "abc123",
		ShortURL:    "http://localhost:8080/s/abc123",
		OriginalURL: "https://example.com",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := `{"short_code":"abc123","short_url":"http://localhost:8080/s/abc123","original_url":"https://example.com"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestAnalytics_EmptySlices(t *testing.T) {
	a := Analytics{
		ShortCode:   "test",
		OriginalURL: "https://example.com",
		TotalClicks: 0,
		ByDay:       []DayStat{},
		ByMonth:     []MonthStat{},
		ByUserAgent: []UserAgentStat{},
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if result["total_clicks"].(float64) != 0 {
		t.Errorf("expected total_clicks 0, got %v", result["total_clicks"])
	}
}
