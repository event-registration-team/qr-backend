CREATE TABLE participants (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    last_name VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    car_number VARCHAR(50),
    qr_token VARCHAR(255) UNIQUE NOT NULL,
    visit_status VARCHAR(50) NOT NULL DEFAULT 'registered' 
        CHECK (visit_status IN ('registered', 'visited')),
    checked_in_at TIMESTAMP WITH TIME ZONE,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_event_email UNIQUE (event_id, email)
);

CREATE TRIGGER set_updated_at_participants 
BEFORE UPDATE ON participants 
FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE INDEX idx_participants_event_id ON participants(event_id);
CREATE INDEX idx_participants_qr_token ON participants(qr_token);
CREATE INDEX idx_participants_email ON participants(email);