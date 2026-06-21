package handlers

import (
	"encoding/json"
	"event-registration/internal/service"
	"event-registration/pkg/utils"
	"net/http"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// Login авторизация администратора
// POST /api/admin/login
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.WriteJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	admin, err := h.service.GetAdminByEmail(loginData.Email)
	if err != nil {
		utils.WriteJSONError(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	if admin == nil {
		utils.WriteJSONError(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// TODO: Здесь нужно проверить хэш пароля
	// if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(loginData.Password)); err != nil {
	// 	http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
	// 	return
	// }

	// TODO: Создать сессию или JWT токен

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Успешная авторизация",
		"admin": map[string]interface{}{
			"id":    admin.ID,
			"email": admin.Email,
		},
	})
}

// GetProfile получение профиля администратора
// GET /api/admin/profile
func (h *AdminHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Получить ID администратора из сессии/токена
	// Здесь заглушка
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Профиль администратора",
	})
}
