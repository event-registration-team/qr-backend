package main

import (
	"event-registration/internal/config"
	"event-registration/internal/database"
	"event-registration/internal/handlers"
	"event-registration/internal/repository"
	"event-registration/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем базу данных
	db := database.Init(cfg.DB)
	defer db.Close()

	// Инициализируем репозитории
	adminRepo := repository.NewAdminRepo(db)
	eventRepo := repository.NewEventRepo(db)
	participantRepo := repository.NewParticipantRepo(db)

	// Инициализируем сервисы
	adminService := service.NewAdminService(adminRepo)
	eventService := service.NewEventService(eventRepo, participantRepo)
	participantService := service.NewParticipantService(participantRepo)

	// Инициализируем хендлеры
	adminHandler := handlers.NewAdminHandler(adminService)
	eventHandler := handlers.NewEventHandler(eventService)
	participantHandler := handlers.NewParticipantHandler(participantService, eventService)

	// Создаем роутер
	r := mux.NewRouter()

	// === API Routes ===

	// Admin routes
	r.HandleFunc("/api/admin/login", adminHandler.Login).Methods("POST")
	r.HandleFunc("/api/admin/profile", adminHandler.GetProfile).Methods("GET")

	// Event routes
	r.HandleFunc("/api/events", eventHandler.CreateEvent).Methods("POST")
	r.HandleFunc("/api/events", eventHandler.GetAllEvents).Methods("GET")
	r.HandleFunc("/api/events/{id}", eventHandler.GetEventByID).Methods("GET")
	r.HandleFunc("/api/events/{id}", eventHandler.UpdateEvent).Methods("PUT")
	r.HandleFunc("/api/events/{id}", eventHandler.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/api/events/{id}/open-registration", eventHandler.OpenRegistration).Methods("POST")
	r.HandleFunc("/api/events/{id}/close-registration", eventHandler.CloseRegistration).Methods("POST")

	// Participant routes
	r.HandleFunc("/api/participants/register", participantHandler.Register).Methods("POST")
	r.HandleFunc("/api/participants/scan", participantHandler.GetByQRToken).Methods("GET")
	r.HandleFunc("/api/participants/{id}/check-in", participantHandler.MarkAsVisited).Methods("POST")
	r.HandleFunc("/api/participants/event", participantHandler.GetByEventID).Methods("GET")
	r.HandleFunc("/api/participants/export", participantHandler.ExportToExcel).Methods("GET")
	r.HandleFunc("/api/participants/import", participantHandler.ImportFromExcel).Methods("POST")
	r.HandleFunc("/api/participants/{id}/qrcode", participantHandler.GetQRCode).Methods("GET")

	// Получаем порт из конфигурации
	port := cfg.Server.Port

	// Запускаем сервер
	fmt.Printf("\n🚀 Сервер запущен на порту %s\n", port)
	fmt.Printf("📍 API доступно по адресу: http://localhost:%s\n", port)
	fmt.Printf("📋 Swagger (если будет): http://localhost:%s/api/docs\n\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}