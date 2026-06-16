package service

import (
	"errors"
	"event-registration/internal/models"
	"event-registration/internal/repository"
)

type EventService struct {
	eventRepo       *repository.EventRepo
	participantRepo *repository.ParticipantRepo
}

func NewEventService(eventRepo *repository.EventRepo, participantRepo *repository.ParticipantRepo) *EventService {
	return &EventService{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
	}
}

// CreateEvent создает новое мероприятие
func (s *EventService) CreateEvent(event *models.Event) error {
	// Валидация обязательных полей
	if event.Title == "" || event.Location == "" {
		return errors.New("название и место проведения обязательны")
	}
	
	// Если статус не указан, устанавливаем по умолчанию
	if event.RegistrationStatus == "" {
		event.RegistrationStatus = "open"
	}

	return s.eventRepo.CreateEvent(event)
}

// GetAllEvents получает список всех мероприятий
func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.eventRepo.GetAllEvents()
}

// GetEventByID получает мероприятие по ID
func (s *EventService) GetEventByID(id int) (*models.Event, error) {
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("мероприятие не найдено")
	}
	return event, nil
}

// UpdateEvent обновляет данные мероприятия
func (s *EventService) UpdateEvent(event *models.Event) error {
	if event.Title == "" || event.Location == "" {
		return errors.New("название и место проведения обязательны")
	}
	return s.eventRepo.UpdateEvent(event)
}

// DeleteEvent удаляет мероприятие
func (s *EventService) DeleteEvent(id int) error {
	return s.eventRepo.DeleteEvent(id)
}

// OpenRegistration открывает регистрацию (меняет статус на 'open')
func (s *EventService) OpenRegistration(id int) error {
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("мероприятие не найдено")
	}
	
	event.RegistrationStatus = "open"
	return s.eventRepo.UpdateEvent(event)
}

// CloseRegistration закрывает регистрацию (меняет статус на 'closed')
func (s *EventService) CloseRegistration(id int) error {
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("мероприятие не найдено")
	}
	
	event.RegistrationStatus = "closed"
	return s.eventRepo.UpdateEvent(event)
}

// CheckRegistrationLimit проверяет, не превышен ли лимит участников
func (s *EventService) CheckRegistrationLimit(eventID int, maxParticipants *int) (bool, error) {
	if maxParticipants == nil {
		return true, nil // Если лимит не установлен, регистрация разрешена
	}
	
	count, err := s.participantRepo.CountParticipantsByEventID(eventID)
	if err != nil {
		return false, err
	}
	
	return count < *maxParticipants, nil
}