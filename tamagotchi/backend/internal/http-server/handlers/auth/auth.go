package handlers

import (
	"encoding/json"
	"net/http"
	"tamagochi-team/siverus/tamagochi/internal/storage/postgres"
)

// LoginRequest — структура запроса для логина.
type LoginRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse — структура ответа на успешную авторизацию.
type LoginResponse struct {
	Message string `json:"message"`
}

// LoginHandler принимает экземпляр хранилища и обрабатывает запрос на авторизацию.
func LoginHandler(storage *postgres.Storage, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Декодирование JSON из тела запроса.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// Вызов метода аутентификации через экземпляр хранилища.
	ok, err := storage.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Успешная авторизация"})
	} else {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
	}
}
