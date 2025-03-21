package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/meetalodariya/email-thread-summarizer/model"
)

const INBOX_PAGE_SIZE = 30

func (h *Handler) getInboxPage(userID string, cursor string, prevCursor string) ([]model.ThreadSummary, string, string, error) {
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to parse cursor: %w", err)
	}
	prevCursorInt, err := strconv.ParseInt(prevCursor, 10, 64)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to parse prevCursor: %w", err)
	}

	query := h.DB.Order("created_at DESC").Limit(INBOX_PAGE_SIZE)

	// Forward pagination (next page)
	if cursorInt > 0 {
		query = query.Where("created_at < ? AND user_id = ?", time.UnixMilli(cursorInt), userID)
	}

	// Backward pagination (previous page)
	if prevCursorInt > 0 {
		query = query.Where("created_at > ? AND user_id = ?", time.UnixMilli(prevCursorInt), userID)
	}

	var threads []model.ThreadSummary
	if err := query.Find(&threads).Error; err != nil {
		return nil, "", "", fmt.Errorf("failed to fetch email summaries: %w", err)
	}

	// Sort threads by creation time descending
	sort.Slice(threads, func(i, j int) bool {
		return threads[i].CreatedAt.After(threads[j].CreatedAt)
	})

	nextCursor := "0"
	prevCursorResult := "0"
	if len(threads) > 0 {
		nextCursor = strconv.FormatInt(threads[len(threads)-1].CreatedAt.UnixMilli(), 10)
		prevCursorResult = strconv.FormatInt(threads[0].CreatedAt.UnixMilli(), 10)
	}

	return threads, nextCursor, prevCursorResult, nil
}

func (h *Handler) HandleGetUserInbox(eCtx echo.Context) error {
	// Extract user ID from JWT token
	user := eCtx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Parse cursor from query params
	cursor, err := decodeCursor(eCtx.QueryParam("cursor"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid cursor")
	}

	// Parse previous cursor from query params
	prevCursor, err := decodeCursor(eCtx.QueryParam("prevCursor"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid prevCursor")
	}

	// Get inbox page
	threads, nextCursor, prevCursorResult, err := h.getInboxPage(userID, cursor, prevCursor)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch email summaries")
	}

	// Encode cursors for response
	nextCursorEnc := b64.StdEncoding.EncodeToString([]byte(nextCursor))
	prevCursorEnc := b64.StdEncoding.EncodeToString([]byte(prevCursorResult))

	return eCtx.JSON(http.StatusOK, JsonResponse{
		Data: threads,
		Pagination: Pagination{
			NextCursor: nextCursorEnc,
			PrevCursor: prevCursorEnc,
		},
	})
}

func decodeCursor(encoded string) (string, error) {
	if encoded == "" {
		return "0", nil
	}

	decoded, err := b64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
