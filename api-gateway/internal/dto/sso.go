package dto

type RegisterRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Role     string `json:"role"`
}

type RegisterResponce struct {
	UserID string `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponce struct {
	Token string `json:"token"`
}

type IsAdminRequest struct {
	ID string `json:"id,omitempty"`
}

type IsAdminResponce struct {
	IsAdmin bool `json:"is_admin"`
}
