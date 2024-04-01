CREATE TABLE IF NOT EXISTS User (
    ID BIGINT PRIMARY KEY AUTO_INCREMENT,
    Email VARCHAR(255) NOT NULL UNIQUE,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    Scope INT,
    IsActive TINYINT DEFAULT 0,
    Version INT
);

CREATE INDEX user_email_idx ON User (Email);


INSERT INTO User (
    ID,
    Email,
    CreatedAt,
    UpdatedAt,
    Scope,
    IsActive,
    Version
) VALUES(
    1,
    "mario@google.com",
    NOW(),
    NOW(),
    0,
    1,
    1
);

