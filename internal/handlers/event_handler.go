package handlers

import (
	"encoding/json"
	"event-registration/internal/models"
	"event-registration/internal/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// CreateEvent создает новое мероприятие
// POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var eventReq struct {
		Title            string  `json:"title"`
		Description      *string `json:"description"`
		Location         string  `json:"location"`
		StartAt          string  `json:"start_at"`
		EndAt            string  `json:"end_at"`
		MaxParticipants  *int    `json:"max_participants"`
		MaterialsLink    *string `json:"materials_link"`
		RequirePhone     bool    `json:"require_phone"`
		RequireCarNumber bool    `json:"require_car_number"`
		RegistrationLink string  `json:"registration_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&eventReq); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Парсим даты
	startAt, err := time.Parse("2006-01-02T15:04:05", eventReq.StartAt)
	if err != nil {
		http.Error(w, "Неверный формат даты начала. Используйте формат: 2006-01-02T15:04:05", http.StatusBadRequest)
		return
	}

	endAt, err := time.Parse("2006-01-02T15:04:05", eventReq.EndAt)
	if err != nil {
		http.Error(w, "Неверный формат даты окончания. Используйте формат: 2006-01-02T15:04:05", http.StatusBadRequest)
		return
	}

	// Генерируем ссылку на регистрацию, если не передана
	if eventReq.RegistrationLink == "" {
		eventReq.RegistrationLink = fmt.Sprintf("event-%d", time.Now().UnixNano())
	}

	// Создаем модель мероприятия
	event := &models.Event{
		Title:              eventReq.Title,
		Description:        eventReq.Description,
		Location:           eventReq.Location,
		StartAt:            startAt,
		EndAt:              endAt,
		RegistrationStatus: "open",
		RegistrationLink:   eventReq.RegistrationLink,
		MaxParticipants:    eventReq.MaxParticipants,
		MaterialsLink:      eventReq.MaterialsLink,
		RequirePhone:       eventReq.RequirePhone,
		RequireCarNumber:   eventReq.RequireCarNumber,
	}

	if err := h.service.CreateEvent(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// GetAllEvents получает список всех мероприятий
// GET /api/events
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// GetEventByID получает мероприятие по ID
// GET /api/events/{id}
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	event, err := h.service.GetEventByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	stats, err := h.service.GetStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                         event.ID,
		"title":                      event.Title,
		"description":                event.Description,
		"location":                   event.Location,
		"start_at":                   event.StartAt,
		"end_at":                     event.EndAt,
		"registration_status":        event.RegistrationStatus,
		"registration_link":          event.RegistrationLink,
		"max_participants":           event.MaxParticipants,
		"materials_link":             event.MaterialsLink,
		"require_phone":              event.RequirePhone,
		"require_car_number":         event.RequireCarNumber,
		"created_at":                 event.CreatedAt,
		"updated_at":                 event.UpdatedAt,
		"current_participants_count": stats["total"],
	})
}

// UpdateEvent обновляет мероприятие (только переданные поля)
// PATCH /api/events/{id}
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Получаем существующее мероприятие
	event, err := h.service.GetEventByID(id)
	if err != nil {
		http.Error(w, "Мероприятие не найдено", http.StatusNotFound)
		return
	}

	// Декодируем только переданные поля через map
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновляем только то что пришло
	if v, ok := req["title"].(string); ok && v != "" {
		event.Title = v
	}
	if v, ok := req["location"].(string); ok && v != "" {
		event.Location = v
	}
	if v, ok := req["description"].(string); ok {
		event.Description = &v
	}
	if v, ok := req["registration_status"].(string); ok && v != "" {
		event.RegistrationStatus = v
	}
	if v, ok := req["registration_link"].(string); ok && v != "" {
		event.RegistrationLink = v
	}
	if v, ok := req["materials_link"].(string); ok {
		event.MaterialsLink = &v
	}
	if v, ok := req["require_phone"].(bool); ok {
		event.RequirePhone = v
	}
	if v, ok := req["require_car_number"].(bool); ok {
		event.RequireCarNumber = v
	}
	if v, ok := req["max_participants"].(float64); ok {
		mp := int(v)
		event.MaxParticipants = &mp
	}
	if v, ok := req["start_at"].(string); ok && v != "" {
		t, err := time.Parse("2006-01-02T15:04:05", v)
		if err != nil {
			http.Error(w, "Неверный формат start_at", http.StatusBadRequest)
			return
		}
		event.StartAt = t
	}
	if v, ok := req["end_at"].(string); ok && v != "" {
		t, err := time.Parse("2006-01-02T15:04:05", v)
		if err != nil {
			http.Error(w, "Неверный формат end_at", http.StatusBadRequest)
			return
		}
		event.EndAt = t
	}

	if err := h.service.UpdateEvent(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DeleteEvent удаляет мероприятие
// DELETE /api/events/{id}
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteEvent(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Мероприятие удалено"})
}

// OpenRegistration открывает регистрацию
// POST /api/events/{id}/open-registration
func (h *EventHandler) OpenRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.service.OpenRegistration(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Регистрация открыта"})
}

// CloseRegistration закрывает регистрацию
// POST /api/events/{id}/close-registration
func (h *EventHandler) CloseRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.service.CloseRegistration(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Регистрация закрыта"})
}

// GetStats возвращает статистику по мероприятию
// GET /api/events/{id}/stats
func (h *EventHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	stats, err := h.service.GetStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
