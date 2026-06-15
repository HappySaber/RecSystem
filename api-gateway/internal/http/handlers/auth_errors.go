package handlers

import (
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func writePlainError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}

func respondLoginError(w http.ResponseWriter, err error) {
	writePlainError(w, http.StatusUnauthorized, "Неверный логин или пароль")
}

func respondRegisterError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.AlreadyExists:
			writePlainError(w, http.StatusConflict, "Пользователь с таким email уже зарегистрирован")
			return
		case codes.InvalidArgument:
			writePlainError(w, http.StatusBadRequest, "Проверьте правильность заполнения полей")
			return
		}
	}

	if strings.Contains(strings.ToLower(err.Error()), "already exists") {
		writePlainError(w, http.StatusConflict, "Пользователь с таким email уже зарегистрирован")
		return
	}

	writePlainError(w, http.StatusInternalServerError, "Не удалось зарегистрироваться. Попробуйте позже")
}
