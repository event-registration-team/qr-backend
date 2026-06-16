package service

import (
	"errors"
	"event-registration/internal/models"
	"event-registration/internal/repository"
)

type AdminService struct {
	repo *repository.AdminRepo
}

func NewAdminService(repo *repository.AdminRepo) *AdminService {
	return &AdminService{repo: repo}
}

// CreateAdmin создает нового администратора
func (s *AdminService) CreateAdmin(admin *models.Admin) error {
	if admin.Email == "" || admin.PasswordHash == "" {
		return errors.New("email и пароль обязательны")
	}
	
	// Проверяем, существует ли админ с таким email
	existing, err := s.repo.GetAdminByEmail(admin.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("администратор с таким email уже существует")
	}
	
	return s.repo.CreateAdmin(admin)
}

// GetAdminByEmail находит администратора по email (для входа)
func (s *AdminService) GetAdminByEmail(email string) (*models.Admin, error) {
	return s.repo.GetAdminByEmail(email)
}

// GetAdminByID находит администратора по ID
func (s *AdminService) GetAdminByID(id int) (*models.Admin, error) {
	admin, err := s.repo.GetAdminByID(id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("администратор не найден")
	}
	return admin, nil
}