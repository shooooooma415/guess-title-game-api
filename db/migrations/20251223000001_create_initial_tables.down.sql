-- Drop tables in reverse order (to respect foreign key constraints)
DROP TABLE IF EXISTS room_emojis;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS users;

-- Drop ENUM types
DROP TYPE IF EXISTS participant_role;
DROP TYPE IF EXISTS room_status;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
