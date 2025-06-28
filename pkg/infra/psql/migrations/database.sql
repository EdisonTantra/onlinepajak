CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE user_customer (
	id BIGSERIAL PRIMARY KEY,
	external_id uuid UNIQUE DEFAULT gen_random_uuid(),
	username VARCHAR(100) UNIQUE NOT NULL,
	password TEXT NOT NULL,
	email VARCHAR(255) UNIQUE  NOT NULL,
	phone_number VARCHAR (20) UNIQUE NOT NULL,
	is_premium BOOLEAN DEFAULT FALSE,
	is_verified BOOLEAN DEFAULT FALSE,
	is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TYPE gender AS ENUM ('MALE', 'FEMALE');

CREATE TABLE profile (
    id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL UNIQUE REFERENCES user_customer(id),
	first_name VARCHAR (60) NOT NULL,
	last_name VARCHAR (60) NOT NULL,
    gender gender NULL,
    age INT NULL,
    description TEXT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NULL
);

-- function for update updated_at
CREATE FUNCTION update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- create triggers
CREATE TRIGGER user_customer_updated_at BEFORE
UPDATE
    ON
    user_customer FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER profile_updated_at BEFORE
UPDATE
    ON
    profile FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- input initial data
INSERT INTO user_customer (id, external_id, username, password, email, phone_number)
VALUES (1, 'ef4a02f6-666e-4233-a890-3f1087dde73d', 'test1', crypt('test1', gen_salt('bf')), 'edisontantra+1@gmail.com', '+62812345677');
INSERT INTO user_customer (id, external_id, username, password, email, phone_number)
VALUES (2, 'c25f862e-bef7-4b71-b9dc-c8c097937a7e', 'test2', crypt('test2', gen_salt('bf')), 'edisontantra+2@gmail.com', '+62812345678');


INSERT INTO profile (id, user_id, first_name, last_name)
VALUES (1, 1, 'edison', 'tantra 1');

INSERT INTO profile (id, user_id, first_name, last_name)
VALUES (2, 2, 'edison', 'tantra 2');