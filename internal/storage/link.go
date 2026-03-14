package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eztwokey/l3-shortener/internal/models"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("short code already exists")
)

func cacheKey(code string) string {
	return "short:" + code
}

const cacheTTL = 24 * time.Hour

func (s *Storage) CreateLink(ctx context.Context, link models.Link) (models.Link, error) {
	var id int
	var createdAt time.Time

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO links (short_code, original_url)
		 VALUES ($1, $2)
		 RETURNING id, created_at`,
		link.ShortCode, link.OriginalURL,
	).Scan(&id, &createdAt)

	if err != nil {
		if isUniqueViolation(err) {
			return models.Link{}, ErrAlreadyExists
		}
		return models.Link{}, fmt.Errorf("insert link: %w", err)
	}

	link.ID = id
	link.CreatedAt = createdAt

	_ = s.rdb.SetWithExpiration(ctx, cacheKey(link.ShortCode), link.OriginalURL, cacheTTL)

	return link, nil
}

func (s *Storage) GetLinkByCode(ctx context.Context, code string) (models.Link, error) {
	cached, err := s.rdb.Get(ctx, cacheKey(code))
	if err == nil && cached != "" {
		return models.Link{ShortCode: code, OriginalURL: cached}, nil
	}

	var link models.Link
	err = s.db.QueryRowContext(ctx,
		`SELECT id, short_code, original_url, created_at
		 FROM links WHERE short_code = $1`,
		code,
	).Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Link{}, ErrNotFound
		}
		return models.Link{}, fmt.Errorf("get link: %w", err)
	}

	_ = s.rdb.SetWithExpiration(ctx, cacheKey(code), link.OriginalURL, cacheTTL)

	return link, nil
}

func (s *Storage) GetLinkFullByCode(ctx context.Context, code string) (models.Link, error) {
	var link models.Link
	err := s.db.QueryRowContext(ctx,
		`SELECT id, short_code, original_url, created_at
		 FROM links WHERE short_code = $1`,
		code,
	).Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Link{}, ErrNotFound
		}
		return models.Link{}, fmt.Errorf("get link full: %w", err)
	}

	return link, nil
}
func isUniqueViolation(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "23505"))
}
