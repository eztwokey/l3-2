package interfaces

import (
	"context"

	"github.com/eztwokey/l3-shortener/internal/models"
)

type LinkStorage interface {
	CreateLink(ctx context.Context, link models.Link) (models.Link, error)
	GetLinkByCode(ctx context.Context, code string) (models.Link, error)
	GetLinkFullByCode(ctx context.Context, code string) (models.Link, error)
}

type ClickStorage interface {
	SaveClick(ctx context.Context, click models.Click) error
	GetTotalClicks(ctx context.Context, linkID int) (int, error)
	GetClicksByDay(ctx context.Context, linkID int) ([]models.DayStat, error)
	GetClicksByMonth(ctx context.Context, linkID int) ([]models.MonthStat, error)
	GetClicksByUserAgent(ctx context.Context, linkID int) ([]models.UserAgentStat, error)
}

type Store interface {
	LinkStorage
	ClickStorage
}
