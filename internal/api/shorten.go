package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eztwokey/l3-shortener/internal/logic"
	"github.com/eztwokey/l3-shortener/internal/models"
	"github.com/eztwokey/l3-shortener/internal/storage"
)

func (a *Api) createLink(c *gin.Context) {
	var req models.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("shorten: bind error", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	link, err := a.logic.CreateLink(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	shortURL := fmt.Sprintf("%s/s/%s", getBaseURL(c), link.ShortCode)

	c.JSON(http.StatusCreated, models.CreateLinkResponse{
		ShortCode:   link.ShortCode,
		ShortURL:    shortURL,
		OriginalURL: link.OriginalURL,
	})
}

func (a *Api) redirect(c *gin.Context) {
	code := c.Param("code")

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	originalURL, err := a.logic.Redirect(c.Request.Context(), code, userAgent, ip)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func (a *Api) getAnalytics(c *gin.Context) {
	code := c.Param("code")

	analytics, err := a.logic.GetAnalytics(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
