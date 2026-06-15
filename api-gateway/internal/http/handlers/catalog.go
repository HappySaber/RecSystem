package handlers

import (
	"api-gateway/internal/dto"
	"api-gateway/internal/events/schemas"
	"api-gateway/internal/middleware"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CatalogHandler struct {
	Producer      EventProducer
	CatalogClient CatalogClient
}

type EventProducer interface {
	Publish(ctx context.Context, topic, key string, event schemas.UserActionEvent) error
}

type CatalogClient interface {
	GetContent(ctx context.Context, contentID string) (*dto.Content, error)
	GetContentByIDs(ctx context.Context, ids []string) ([]dto.ContentShort, error)
	FindContentByExternal(ctx context.Context, contentID, externalSource string) (*dto.Content, error)
	GetMovieDetails(ctx context.Context, contentID string) (*dto.MovieDetails, error)
}

func (h *CatalogHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	contentID := parseContentID(r)
	if contentID == "" {
		http.Error(w, "content_id is required", http.StatusBadRequest)
		return
	}

	resp, err := h.CatalogClient.GetContent(r.Context(), contentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()

	fmt.Println(schemas.UserActionEvent{
		UserID:    userID,
		ContentID: contentID,
		Action:    "VIEW",
		Timestamp: time.Now()})
	err = h.Producer.Publish(ctx, "user-action", contentID, schemas.UserActionEvent{
		UserID:    userID,
		ContentID: contentID,
		Action:    "VIEW",
		Timestamp: time.Now(),
	})
	if err != nil {
		fmt.Println("failed to publish user action event:", err)
	}
}

func (h *CatalogHandler) GetContentDetails(w http.ResponseWriter, r *http.Request) {
	contentID := mux.Vars(r)["content_id"]
	if contentID == "" {
		http.Error(w, "content_id is required", http.StatusBadRequest)
		return
	}

	content, err := h.CatalogClient.GetContent(r.Context(), contentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ContentDetailsResponse{Content: *content}

	if content.Type == "movie" {
		if movie, err := h.CatalogClient.GetMovieDetails(r.Context(), contentID); err == nil {
			resp.Movie = movie
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		return
	}

	ctx := context.Background()
	_ = h.Producer.Publish(ctx, "user-action", contentID, schemas.UserActionEvent{
		UserID:    userID,
		ContentID: contentID,
		Action:    "VIEW",
		Timestamp: time.Now(),
	})
}

func (h *CatalogHandler) GetContentByIDs(w http.ResponseWriter, r *http.Request) {
	var req dto.GetContentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.CatalogClient.GetContent(r.Context(), req.ContentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) FindContentByExternal(w http.ResponseWriter, r *http.Request) {
	var req dto.FindContentByExternalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.CatalogClient.FindContentByExternal(r.Context(), req.ContentID, req.ExternalSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
