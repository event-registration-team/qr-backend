package repository

import (
	"database/sql"
	"event-registration/internal/models"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

// CreateAdmin создает нового администратора
func (r *AdminRepo) CreateAdmin(admin *models.Admin) error {
	query := `
		INSERT INTO admins (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		admin.Email,
		admin.PasswordHash,
	).Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt)

	return err
}

// GetAdminByEmail находит администратора по email (для входа)
func (r *AdminRepo) GetAdminByEmail(email string) (*models.Admin, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM admins WHERE email = $1`
	
	var admin models.Admin
	err := r.db.QueryRow(query, email).Scan(
		&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Админ не найден
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetAdminByID находит администратора по ID
func (r *AdminRepo) GetAdminByID(id int) (*models.Admin, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM admins WHERE id = $1`
	
	var admin models.Admin
	err := r.db.QueryRow(query, id).Scan(
		&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}