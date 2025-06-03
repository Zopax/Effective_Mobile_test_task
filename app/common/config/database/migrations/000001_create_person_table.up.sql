CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    surname VARCHAR(150) NOT NULL,
    patronymic VARCHAR(150),
    age INT NOT NULL,
    gender VARCHAR(10) NOT NULL,
    nationality VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_persons_name ON persons(name);
CREATE INDEX IF NOT EXISTS idx_persons_surname ON persons(surname);
CREATE INDEX IF NOT EXISTS idx_persons_age ON persons(age);
CREATE INDEX IF NOT EXISTS idx_persons_gender ON persons(gender);
CREATE INDEX IF NOT EXISTS idx_persons_nationality ON persons(nationality);