package storage

import (
	"context"
	"fmt"

	"github.com/eztwokey/l3-shortener/internal/models"
)

func (s *Storage) SaveClick(ctx context.Context, click models.Click) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO clicks (link_id, user_agent, ip_address)
		 VALUES ($1, $2, $3)`,
		click.LinkID, click.UserAgent, click.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("save click: %w", err)
	}
	return nil
}

func (s *Storage) GetTotalClicks(ctx context.Context, linkID int) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM clicks WHERE link_id = $1`,
		linkID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get total clicks: %w", err)
	}
	return count, nil
}

func (s *Storage) GetClicksByDay(ctx context.Context, linkID int) ([]models.DayStat, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT TO_CHAR(clicked_at, 'YYYY-MM-DD') AS day, COUNT(*) AS clicks
		 FROM clicks
		 WHERE link_id = $1
		 GROUP BY day
		 ORDER BY day DESC`,
		linkID,
	)
	if err != nil {
		return nil, fmt.Errorf("get clicks by day: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []models.DayStat
	for rows.Next() {
		var s models.DayStat
		if err := rows.Scan(&s.Date, &s.Clicks); err != nil {
			return nil, fmt.Errorf("scan day stat: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (s *Storage) GetClicksByMonth(ctx context.Context, linkID int) ([]models.MonthStat, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT TO_CHAR(clicked_at, 'YYYY-MM') AS month, COUNT(*) AS clicks
		 FROM clicks
		 WHERE link_id = $1
		 GROUP BY month
		 ORDER BY month DESC`,
		linkID,
	)
	if err != nil {
		return nil, fmt.Errorf("get clicks by month: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []models.MonthStat
	for rows.Next() {
		var s models.MonthStat
		if err := rows.Scan(&s.Month, &s.Clicks); err != nil {
			return nil, fmt.Errorf("scan month stat: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (s *Storage) GetClicksByUserAgent(ctx context.Context, linkID int) ([]models.UserAgentStat, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT user_agent, COUNT(*) AS clicks
		 FROM clicks
		 WHERE link_id = $1
		 GROUP BY user_agent
		 ORDER BY clicks DESC`,
		linkID,
	)
	if err != nil {
		return nil, fmt.Errorf("get clicks by user agent: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []models.UserAgentStat
	for rows.Next() {
		var s models.UserAgentStat
		if err := rows.Scan(&s.UserAgent, &s.Clicks); err != nil {
			return nil, fmt.Errorf("scan user agent stat: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}
