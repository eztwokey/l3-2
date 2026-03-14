package logic

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/eztwokey/l3-shortener/internal/models"
	"github.com/eztwokey/l3-shortener/internal/shortgen"
	"github.com/eztwokey/l3-shortener/internal/storage"
)

var (
	ErrBadRequest = errors.New("bad request")
)

const maxGenerateAttempts = 5

func (l *Logic) CreateLink(ctx context.Context, req models.CreateLinkRequest) (models.Link, error) {
	rawURL := strings.TrimSpace(req.URL)
	if rawURL == "" {
		return models.Link{}, ErrBadRequest
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return models.Link{}, fmt.Errorf("%w: invalid url", ErrBadRequest)
	}

	code := strings.TrimSpace(req.CustomCode)

	if code != "" {
		// Кастомный код — проверяем длину и символы
		if len(code) < 3 || len(code) > 20 {
			return models.Link{}, fmt.Errorf("%w: custom code must be 3-20 characters", ErrBadRequest)
		}
	} else {
		for i := 0; i < maxGenerateAttempts; i++ {
			generated, err := shortgen.Generate()
			if err != nil {
				return models.Link{}, fmt.Errorf("generate code: %w", err)
			}
			code = generated

			link, err := l.store.CreateLink(ctx, models.Link{
				ShortCode:   code,
				OriginalURL: rawURL,
			})
			if err != nil {
				if errors.Is(err, storage.ErrAlreadyExists) {
					l.logger.Warn("short code collision, retrying",
						"code", code,
						"attempt", i+1,
					)
					continue
				}
				l.logger.Error("create link failed", "err", err)
				return models.Link{}, err
			}

			l.logger.Info("link created", "code", link.ShortCode, "url", link.OriginalURL)
			return link, nil
		}

		return models.Link{}, fmt.Errorf("failed to generate unique code after %d attempts", maxGenerateAttempts)
	}

	link, err := l.store.CreateLink(ctx, models.Link{
		ShortCode:   code,
		OriginalURL: rawURL,
	})
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return models.Link{}, fmt.Errorf("%w: code '%s' already taken", ErrBadRequest, code)
		}
		l.logger.Error("create link failed", "err", err)
		return models.Link{}, err
	}

	l.logger.Info("link created", "code", link.ShortCode, "url", link.OriginalURL)
	return link, nil
}

func (l *Logic) Redirect(ctx context.Context, code string, userAgent string, ip string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", ErrBadRequest
	}

	link, err := l.store.GetLinkByCode(ctx, code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", storage.ErrNotFound
		}
		l.logger.Error("redirect: get link failed", "code", code, "err", err)
		return "", err
	}

	go func() {
		bgCtx := context.Background()

		linkID := link.ID
		if linkID == 0 {
			full, err := l.store.GetLinkFullByCode(bgCtx, code)
			if err != nil {
				l.logger.Error("redirect: get full link for click failed", "code", code, "err", err)
				return
			}
			linkID = full.ID
		}

		click := models.Click{
			LinkID:    linkID,
			UserAgent: userAgent,
			IPAddress: ip,
		}
		if err := l.store.SaveClick(bgCtx, click); err != nil {
			l.logger.Error("redirect: save click failed", "code", code, "err", err)
		}
	}()

	return link.OriginalURL, nil
}

func (l *Logic) GetAnalytics(ctx context.Context, code string) (models.Analytics, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return models.Analytics{}, ErrBadRequest
	}

	link, err := l.store.GetLinkFullByCode(ctx, code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Analytics{}, storage.ErrNotFound
		}
		l.logger.Error("analytics: get link failed", "code", code, "err", err)
		return models.Analytics{}, err
	}

	total, err := l.store.GetTotalClicks(ctx, link.ID)
	if err != nil {
		l.logger.Error("analytics: get total clicks failed", "code", code, "err", err)
		return models.Analytics{}, err
	}

	byDay, err := l.store.GetClicksByDay(ctx, link.ID)
	if err != nil {
		l.logger.Error("analytics: get clicks by day failed", "code", code, "err", err)
		return models.Analytics{}, err
	}

	byMonth, err := l.store.GetClicksByMonth(ctx, link.ID)
	if err != nil {
		l.logger.Error("analytics: get clicks by month failed", "code", code, "err", err)
		return models.Analytics{}, err
	}

	byUA, err := l.store.GetClicksByUserAgent(ctx, link.ID)
	if err != nil {
		l.logger.Error("analytics: get clicks by user agent failed", "code", code, "err", err)
		return models.Analytics{}, err
	}

	return models.Analytics{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		TotalClicks: total,
		ByDay:       byDay,
		ByMonth:     byMonth,
		ByUserAgent: byUA,
	}, nil
}
