package service

import (
	"errors"
	"event-registration/internal/models"
	"event-registration/internal/repository"

	"github.com/google/uuid"
)

type ParticipantService struct {
	repo *repository.ParticipantRepo
}

func NewParticipantService(repo *repository.ParticipantRepo) *ParticipantService {
	return &ParticipantService{repo: repo}
}

// CreateParticipant создает нового участника
func (s *ParticipantService) CreateParticipant(participant *models.Participant, maxParticipants *int) error {
	if participant.LastName == "" || participant.FirstName == "" || participant.Email == "" {
		return errors.New("фамилия, имя и email обязательны")
	}

	existing, err := s.repo.GetParticipantByEventIDAndEmail(participant.EventID, participant.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("участник с таким email уже зарегистрирован на это мероприятие")
	}

	if participant.QRToken == "" {
		participant.QRToken = generateQRToken()
	}

	participant.VisitStatus = "registered"

	return s.repo.CreateParticipantTx(participant, maxParticipants)
}

// GetParticipantByID получает участника по ID
func (s *ParticipantService) GetParticipantByID(id int) (*models.Participant, error) {
	participant, err := s.repo.GetParticipantByID(id)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("участник не найден")
	}
	return participant, nil
}

// GetParticipantByQRToken находит участника по QR-токену
func (s *ParticipantService) GetParticipantByQRToken(qrToken string) (*models.Participant, error) {
	participant, err := s.repo.GetParticipantByQRToken(qrToken)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("участник не найден")
	}
	return participant, nil
}

// GetParticipantsByEventID получает всех участников мероприятия
func (s *ParticipantService) GetParticipantsByEventID(eventID int) ([]models.Participant, error) {
	return s.repo.GetParticipantsByEventID(eventID)
}

// MarkAsVisited отмечает участника как посетившего мероприятие
func (s *ParticipantService) MarkAsVisited(id int) error {
	participant, err := s.repo.GetParticipantByID(id)
	if err != nil {
		return errors.New("участник не найден")
	}

	if participant.VisitStatus == "visited" {
		return errors.New("участник уже отмечен как посетивший мероприятие")
	}

	return s.repo.MarkAsVisited(id)
}

// UpdateParticipant обновляет данные участника
func (s *ParticipantService) UpdateParticipant(participant *models.Participant) error {
	return s.repo.UpdateParticipant(participant)
}

// DeleteParticipant удаляет участника
func (s *ParticipantService) DeleteParticipant(id int) error {
	return s.repo.DeleteParticipant(id)
}

// CountParticipantsByEventID считает количество участников
func (s *ParticipantService) CountParticipantsByEventID(eventID int) (int, error) {
	return s.repo.CountParticipantsByEventID(eventID)
}

// generateQRToken генерирует уникальный QR-токен с использованием UUID
func generateQRToken() string {
	return uuid.New().String()
}
