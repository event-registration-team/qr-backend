CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255) NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE NOT NULL,
    registration_status VARCHAR(50) NOT NULL DEFAULT 'open' 
        CHECK (registration_status IN ('open', 'closed', 'completed')),
    registration_link VARCHAR(255) UNIQUE NOT NULL,
    max_participants INT,
    materials_link VARCHAR(255),
    require_phone BOOLEAN DEFAULT FALSE,
    require_car_number BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_updated_at_events 
BEFORE UPDATE ON events 
FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();