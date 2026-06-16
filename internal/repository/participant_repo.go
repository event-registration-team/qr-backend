package repository

import (
	"database/sql"
	"event-registration/internal/models"
	"time"
)

type ParticipantRepo struct {
	db *sql.DB
}

func NewParticipantRepo(db *sql.DB) *ParticipantRepo {
	return &ParticipantRepo{db: db}
}

// CreateParticipant создает нового участника
func (r *ParticipantRepo) CreateParticipant(participant *models.Participant) error {
	query := `
		INSERT INTO participants 
		(event_id, last_name, first_name, middle_name, email, phone, car_number, qr_token, visit_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, registered_at, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		participant.EventID,
		participant.LastName,
		participant.FirstName,
		participant.MiddleName,
		participant.Email,
		participant.Phone,
		participant.CarNumber,
		participant.QRToken,
		participant.VisitStatus,
	).Scan(&participant.ID, &participant.RegisteredAt, &participant.CreatedAt, &participant.UpdatedAt)

	return err
}

// GetParticipantByID получает участника по ID
func (r *ParticipantRepo) GetParticipantByID(id int) (*models.Participant, error) {
	query := `
		SELECT id, event_id, last_name, first_name, middle_name, email, phone, car_number, 
		       qr_token, visit_status, checked_in_at, registered_at, created_at, updated_at 
		FROM participants WHERE id = $1`
	
	var p models.Participant
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.EventID, &p.LastName, &p.FirstName, &p.MiddleName, &p.Email,
		&p.Phone, &p.CarNumber, &p.QRToken, &p.VisitStatus, &p.CheckedInAt,
		&p.RegisteredAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetParticipantByQRToken находит участника по QR-токену (для сканирования)
func (r *ParticipantRepo) GetParticipantByQRToken(qrToken string) (*models.Participant, error) {
	query := `
		SELECT id, event_id, last_name, first_name, middle_name, email, phone, car_number, 
		       qr_token, visit_status, checked_in_at, registered_at, created_at, updated_at 
		FROM participants WHERE qr_token = $1`
	
	var p models.Participant
	err := r.db.QueryRow(query, qrToken).Scan(
		&p.ID, &p.EventID, &p.LastName, &p.FirstName, &p.MiddleName, &p.Email,
		&p.Phone, &p.CarNumber, &p.QRToken, &p.VisitStatus, &p.CheckedInAt,
		&p.RegisteredAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetParticipantsByEventID получает всех участников мероприятия
func (r *ParticipantRepo) GetParticipantsByEventID(eventID int) ([]models.Participant, error) {
	query := `
		SELECT id, event_id, last_name, first_name, middle_name, email, phone, car_number, 
		       qr_token, visit_status, checked_in_at, registered_at, created_at, updated_at 
		FROM participants WHERE event_id = $1 ORDER BY registered_at DESC`
	
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		err := rows.Scan(
			&p.ID, &p.EventID, &p.LastName, &p.FirstName, &p.MiddleName, &p.Email,
			&p.Phone, &p.CarNumber, &p.QRToken, &p.VisitStatus, &p.CheckedInAt,
			&p.RegisteredAt, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

// UpdateParticipant обновляет данные участника
func (r *ParticipantRepo) UpdateParticipant(participant *models.Participant) error {
	query := `
		UPDATE participants 
		SET last_name = $1, first_name = $2, middle_name = $3, email = $4, 
		    phone = $5, car_number = $6, visit_status = $7, checked_in_at = $8, updated_at = $9
		WHERE id = $10`

	result, err := r.db.Exec(
		query,
		participant.LastName, participant.FirstName, participant.MiddleName, participant.Email,
		participant.Phone, participant.CarNumber, participant.VisitStatus, participant.CheckedInAt,
		time.Now(), participant.ID,
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

// MarkAsVisited отмечает участника как посетившего мероприятие (при сканировании QR)
func (r *ParticipantRepo) MarkAsVisited(id int) error {
	query := `
		UPDATE participants 
		SET visit_status = 'visited', checked_in_at = $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.Exec(query, time.Now(), time.Now(), id)
	return err
}

// CountParticipantsByEventID считает количество участников мероприятия (для проверки лимита)
func (r *ParticipantRepo) CountParticipantsByEventID(eventID int) (int, error) {
	query := `SELECT COUNT(*) FROM participants WHERE event_id = $1`
	
	var count int
	err := r.db.QueryRow(query, eventID).Scan(&count)
	return count, err
}

// DeleteParticipant удаляет участника
func (r *ParticipantRepo) DeleteParticipant(id int) error {
	query := `DELETE FROM participants WHERE id = $1`
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
// GetParticipantByEventIDAndEmail находит участника по event_id и email
func (r *ParticipantRepo) GetParticipantByEventIDAndEmail(eventID int, email string) (*models.Participant, error) {
	query := `
		SELECT id, event_id, last_name, first_name, middle_name, email, phone, car_number, 
		       qr_token, visit_status, checked_in_at, registered_at, created_at, updated_at 
		FROM participants WHERE event_id = $1 AND email = $2`
	
	var p models.Participant
	err := r.db.QueryRow(query, eventID, email).Scan(
		&p.ID, &p.EventID, &p.LastName, &p.FirstName, &p.MiddleName, &p.Email,
		&p.Phone, &p.CarNumber, &p.QRToken, &p.VisitStatus, &p.CheckedInAt,
		&p.RegisteredAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}