package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

const defaultLimit = 12

func parseLimit(r *http.Request) int {
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			return n
		}
	}

	body, err := io.ReadAll(r.Body)
	if err == nil && len(body) > 0 {
		var req struct {
			Limit int `json:"limit"`
		}
		if json.Unmarshal(body, &req) == nil && req.Limit > 0 {
			return req.Limit
		}
	}

	return defaultLimit
}

func parseOffset(r *http.Request) int {
	if q := r.URL.Query().Get("offset"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n >= 0 {
			return n
		}
	}
	return 0
}

func parseContentID(r *http.Request) string {
	if id := r.URL.Query().Get("content_id"); id != "" {
		return id
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		return ""
	}

	var req struct {
		ContentID string `json:"content_id"`
	}
	if json.Unmarshal(body, &req) != nil {
		return ""
	}

	return req.ContentID
}

func sliceFromOffset(items []string, offset, limit int) []string {
	if offset >= len(items) {
		return nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}
