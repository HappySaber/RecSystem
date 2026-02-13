package handlers

import (
	"api-gateway/internal/dto"
	"context"
	"encoding/json"
	"net/http"
)

type SSOHandler struct {
	Producer  EventProducer
	SSOClient SSOClient
}

type SSOClient interface {
	Register(ctx context.Context, email, password, name, surname, role string) (*dto.RegisterResponce, error)
	Login(ctx context.Context, email, password string) (*dto.LoginResponce, error)
	IsAdmin(ctx context.Context, email string) (bool, error)
}

func (h *SSOHandler) Register(w *http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	json.NewDecoder(r.Body).Decode(&req)

	if req.Role == "" {
		req.Role = "user"
	}
	resp, err := h.SSOClient.Register(r.Context(), req.Email, req.Password, req.Name, req.Surname, req.Role)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(*w).Encode(resp)
}

func (h *SSOHandler) Login(w *http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.SSOClient.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(*w).Encode(resp)
}

func (h *SSOHandler) IsAdmin(w *http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.SSOClient.IsAdmin(r.Context(), req.Email)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(*w).Encode(resp)
}
