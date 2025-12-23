-- Create ENUM types
CREATE TYPE room_status AS ENUM ('waiting', 'setting_topic', 'discussing', 'answering', 'checking', 'finished');
CREATE TYPE participant_role AS ENUM ('host', 'player');

-- Create User table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Theme table
CREATE TABLE themes (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    hint TEXT
);

-- Create Room table
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    theme_id UUID NOT NULL,
    topic TEXT,
    answer TEXT,
    status room_status NOT NULL DEFAULT 'waiting',
    host_user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    CONSTRAINT fk_room_theme FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE RESTRICT,
    CONSTRAINT fk_room_host FOREIGN KEY (host_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create Participant table
CREATE TABLE participants (
    id UUID NOT NULL,
    room_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role participant_role NOT NULL,
    is_leader BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id),
    CONSTRAINT fk_participant_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_participant_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create Room_Emoji table
CREATE TABLE room_emojis (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL,
    participant_id UUID,
    emoji VARCHAR(50) NOT NULL,
    CONSTRAINT fk_room_emoji_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_rooms_code ON rooms(code);
CREATE INDEX idx_rooms_host_user_id ON rooms(host_user_id);
CREATE INDEX idx_rooms_theme_id ON rooms(theme_id);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_participants_user_id ON participants(user_id);
CREATE INDEX idx_participants_room_id ON participants(room_id);
CREATE INDEX idx_room_emojis_room_id ON room_emojis(room_id);
