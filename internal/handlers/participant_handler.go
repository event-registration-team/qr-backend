package handlers

import (
	"encoding/json"
	"event-registration/internal/models"
	"event-registration/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

type ParticipantHandler struct {
	service      *service.ParticipantService
	eventService *service.EventService
	emailService *service.EmailService // добавить
}

func NewParticipantHandler(
	service *service.ParticipantService,
	eventService *service.EventService,
	emailService *service.EmailService, // добавить
) *ParticipantHandler {
	return &ParticipantHandler{
		service:      service,
		eventService: eventService,
		emailService: emailService, // добавить
	}
}

// Register регистрирует нового участника
// POST /api/participants/register
func (h *ParticipantHandler) Register(w http.ResponseWriter, r *http.Request) {
	var participant struct {
		EventID    int     `json:"event_id"`
		LastName   string  `json:"last_name"`
		FirstName  string  `json:"first_name"`
		MiddleName *string `json:"middle_name"`
		Email      string  `json:"email"`
		Phone      *string `json:"phone"`
		CarNumber  *string `json:"car_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли мероприятие и открыта ли регистрация
	event, err := h.eventService.GetEventByID(participant.EventID)
	if err != nil {
		http.Error(w, "Мероприятие не найдено", http.StatusNotFound)
		return
	}

	if event.RegistrationStatus != "open" {
		http.Error(w, "Регистрация на это мероприятие закрыта", http.StatusForbidden)
		return
	}

	// Проверяем лимит участников
	allowed, err := h.eventService.CheckRegistrationLimit(event.ID, event.MaxParticipants)
	if err != nil {
		http.Error(w, "Ошибка проверки лимита", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "Извините, регистрация завершена", http.StatusForbidden)
		return
	}

	// Создаем модель участника
	newParticipant := &models.Participant{
		EventID:     participant.EventID,
		LastName:    participant.LastName,
		FirstName:   participant.FirstName,
		MiddleName:  participant.MiddleName,
		Email:       participant.Email,
		Phone:       participant.Phone,
		CarNumber:   participant.CarNumber,
		VisitStatus: "registered",
		QRToken:     "", // Сгенерируется в сервисе автоматически
	}

	// Создаем участника через сервис
	if err := h.service.CreateParticipant(newParticipant, event.MaxParticipants); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go h.emailService.SendRegistrationEmail(newParticipant, event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Регистрация успешна",
		"qr_token":    newParticipant.QRToken,
		"participant": newParticipant,
	})
}

// GetByQRToken получает участника по QR-токену (для сканирования)
// GET /api/participants/scan?token=xxx
func (h *ParticipantHandler) GetByQRToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "QR-токен не указан", http.StatusBadRequest)
		return
	}

	participant, err := h.service.GetParticipantByQRToken(token)
	if err != nil {
		http.Error(w, "Участник не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participant)
}

// MarkAsVisited отмечает участника как посетившего мероприятие
// POST /api/participants/{id}/check-in
func (h *ParticipantHandler) MarkAsVisited(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkAsVisited(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Участник отмечен как посетивший мероприятие",
	})
}

// GetByEventID получает всех участников мероприятия
// GET /api/participants/event?event_id=1
func (h *ParticipantHandler) GetByEventID(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Неверный ID мероприятия", http.StatusBadRequest)
		return
	}

	participants, err := h.service.GetParticipantsByEventID(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}

// ExportToExcel экспортирует участников мероприятия в Excel
// GET /api/participants/export?event_id=1
func (h *ParticipantHandler) ExportToExcel(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Неверный ID мероприятия", http.StatusBadRequest)
		return
	}

	participants, err := h.service.GetParticipantsByEventID(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем Excel файл
	f := excelize.NewFile()
	sheetName := "Участники"

	// Заголовки
	headers := []string{"ID", "Фамилия", "Имя", "Отчество", "Email", "Телефон", "Номер авто", "Дата регистрации", "Статус посещения", "Время входа"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Данные
	for i, p := range participants {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), p.LastName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), p.FirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), p.MiddleName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), p.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), p.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), p.CarNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), p.RegisteredAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), p.VisitStatus)
		if p.CheckedInAt != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), p.CheckedInAt.Format("2006-01-02 15:04:05"))
		}
	}

	// Отправляем файл
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=participants_%d.xlsx", eventID))

	if err := f.Write(w); err != nil {
		http.Error(w, "Ошибка экспорта", http.StatusInternalServerError)
		return
	}
}

// ImportFromExcel импортирует участников из Excel
// POST /api/participants/import?event_id=1
func (h *ParticipantHandler) ImportFromExcel(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Неверный ID мероприятия", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли мероприятие
	event, err := h.eventService.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Мероприятие не найдено", http.StatusNotFound)
		return
	}

	// Получаем файл
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Файл не загружен", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем Excel
	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusBadRequest)
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		http.Error(w, "Ошибка чтения листа", http.StatusBadRequest)
		return
	}

	// Пропускаем заголовок (первая строка), обрабатываем данные
	importedCount := 0
	for i, row := range rows {
		if i == 0 {
			continue // Пропускаем заголовок
		}

		if len(row) < 2 {
			continue // Пропускаем строки с недостаточным количеством данных
		}

		// Создаем участника из Excel
		// row[0] - Фамилия, row[1] - Имя, row[2] - Отчество (опционально), row[3] - Email, и т.д.
		newParticipant := &models.Participant{
			EventID:     eventID,
			LastName:    row[0],
			FirstName:   row[1],
			VisitStatus: "registered",
		}

		if len(row) > 2 && row[2] != "" {
			newParticipant.MiddleName = &row[2]
		}
		if len(row) > 3 && row[3] != "" {
			newParticipant.Email = row[3]
		} else {
			continue // Email обязателен
		}
		if len(row) > 4 && row[4] != "" {
			newParticipant.Phone = &row[4]
		}
		if len(row) > 5 && row[5] != "" {
			newParticipant.CarNumber = &row[5]
		}

		// Создаем участника через сервис
		if err := h.service.CreateParticipant(newParticipant, event.MaxParticipants); err != nil {
			// Пропускаем ошибки (например, дубликат email) и продолжаем импорт
			continue
		}
		importedCount++
	}

	go func(e *models.Event, eID int) {
		participants, _ := h.service.GetParticipantsByEventID(eID)
		h.emailService.SendBulkRegistrationEmails(participants, e)
	}(event, eventID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Импортировано %d участников", importedCount),
	})
}

// GetQRCode генерирует и возвращает QR-код для участника
// GET /api/participants/{id}/qrcode
func (h *ParticipantHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	participant, err := h.service.GetParticipantByID(id)
	if err != nil {
		http.Error(w, "Участник не найден", http.StatusNotFound)
		return
	}

	// Генерируем QR-код
	pngData, err := qrcode.Encode(participant.QRToken, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Ошибка генерации QR-кода", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=qr_%d.png", id))
	w.Write(pngData)
}
