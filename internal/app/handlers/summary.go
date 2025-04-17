package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/meetalodariya/email-thread-summarizer/internal/auth"
	"github.com/meetalodariya/email-thread-summarizer/model"
)

const INBOX_PAGE_SIZE = 10

func (h *Handler) getInboxPage(userID string, cursor string, searchQuery string) ([]model.ThreadSummary, string, bool, error) {
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		return nil, "", false, fmt.Errorf("failed to parse cursor: %w", err)
	}

	query := h.DB
	if cursorInt > 0 {
		query = query.Where("most_recent_email_timestamp < ?", time.UnixMilli(cursorInt))
	}

	if searchQuery != "" {
		query = query.Where(
			"search_vector @@ plainto_tsquery('english', ?)", searchQuery,
		).Select(
			[]string{"*", fmt.Sprintf("ts_rank(search_vector, plainto_tsquery('english', '%s')) AS rank", searchQuery)},
		).Order("rank DESC")
	} else {
		query = query.Order("most_recent_email_timestamp DESC")
	}

	var threads []model.ThreadSummary
	if err := query.Where("user_id = ?", userID).Limit(INBOX_PAGE_SIZE + 1).Find(&threads).Error; err != nil {
		return nil, "", false, fmt.Errorf("failed to fetch email summaries: %w", err)
	}

	nextCursor := ""
	hasNextPage := false
	threadsLength := len(threads)

	if threadsLength == INBOX_PAGE_SIZE+1 {
		threads = threads[:len(threads)-1]
		nextCursor = strconv.FormatInt(threads[len(threads)-1].MostRecentEmailTimestamp.UnixMilli(), 10)
		hasNextPage = true
	}

	return threads, nextCursor, hasNextPage, nil
}

func (h *Handler) HandleGetUserInbox(eCtx echo.Context) error {
	// Extract user ID from JWT token
	user := eCtx.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userID := strconv.FormatUint(uint64(claims.UserID), 10)

	cursor, err := decodeCursor(eCtx.QueryParam("nextCursor"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid cursor")
	}

	searchQuery := eCtx.QueryParam("q")

	// Get inbox page
	threads, nextCursor, hasNextPage, err := h.getInboxPage(userID, cursor, searchQuery)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch email summaries")
	}

	// Encode cursors for response
	nextCursorEnc := b64.StdEncoding.EncodeToString([]byte(nextCursor))

	return eCtx.JSON(http.StatusOK, PaginatedResponse{
		Data: threads,
		Pagination: Pagination{
			NextCursor:  nextCursorEnc,
			HasNextPage: hasNextPage,
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
