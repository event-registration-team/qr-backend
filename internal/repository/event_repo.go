package repository

import (
	"database/sql"
	"event-registration/internal/models"
	"time"
)

type EventRepo struct {
	db *sql.DB
}

// NewEventRepo создает новый экземпляр репозитория
func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db: db}
}

// CreateEvent создает новое мероприятие в БД
func (r *EventRepo) CreateEvent(event *models.Event) error {
	query := `
		INSERT INTO events 
		(title, description, location, start_at, end_at, registration_status, registration_link, max_participants, materials_link, require_phone, require_car_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.Location,
		event.StartAt,
		event.EndAt,
		event.RegistrationStatus,
		event.RegistrationLink,
		event.MaxParticipants,
		event.MaterialsLink,
		event.RequirePhone,
		event.RequireCarNumber,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	return err
}

// GetAllEvents получает список всех мероприятий
func (r *EventRepo) GetAllEvents() ([]models.Event, error) {
	query := `SELECT id, title, description, location, start_at, end_at, registration_status, registration_link, max_participants, materials_link, require_phone, require_car_number, created_at, updated_at FROM events ORDER BY start_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Location, &e.StartAt, &e.EndAt, 
			&e.RegistrationStatus, &e.RegistrationLink, &e.MaxParticipants, &e.MaterialsLink, 
			&e.RequirePhone, &e.RequireCarNumber, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// GetEventByID получает одно мероприятие по ID
func (r *EventRepo) GetEventByID(id int) (*models.Event, error) {
	query := `SELECT id, title, description, location, start_at, end_at, registration_status, registration_link, max_participants, materials_link, require_phone, require_car_number, created_at, updated_at FROM events WHERE id = $1`
	
	var e models.Event
	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.Location, &e.StartAt, &e.EndAt, 
		&e.RegistrationStatus, &e.RegistrationLink, &e.MaxParticipants, &e.MaterialsLink, 
		&e.RequirePhone, &e.RequireCarNumber, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Мероприятие не найдено
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateEvent обновляет данные мероприятия
func (r *EventRepo) UpdateEvent(event *models.Event) error {
	query := `
		UPDATE events 
		SET title = $1, description = $2, location = $3, start_at = $4, end_at = $5, 
		    registration_status = $6, registration_link = $7, max_participants = $8, 
		    materials_link = $9, require_phone = $10, require_car_number = $11, updated_at = $12
		WHERE id = $13`

	result, err := r.db.Exec(
		query,
		event.Title, event.Description, event.Location, event.StartAt, event.EndAt,
		event.RegistrationStatus, event.RegistrationLink, event.MaxParticipants,
		event.MaterialsLink, event.RequirePhone, event.RequireCarNumber, time.Now(), event.ID,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteEvent удаляет мероприятие по ID
func (r *EventRepo) DeleteEvent(id int) error {
	query := `DELETE FROM events WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}