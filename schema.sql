CREATE DATABASE doodle CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE doodle;

CREATE TABLE events (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    time DATETIME NOT NULL
);

CREATE INDEX idx_events_time ON events(time);

INSERT INTO events (title, description, time) VALUES (
    'Festival music',
    'Please make sure to be available all day',
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

INSERT INTO events (title, description, time) VALUES (
    'Jam session',
    'Will be so cool',
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 17 DAY)
);

INSERT INTO events (title, description, time) VALUES (
    'Concert of classical music',
    'Suit yourself!',
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 27 DAY)
);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

CREATE USER 'web'@'%';
GRANT SELECT, INSERT, UPDATE ON doodle.* TO 'web'@'%';
ALTER USER 'web'@'%' IDENTIFIED BY 'pass';
