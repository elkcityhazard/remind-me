CREATE TABLE IF NOT EXISTS ActivationToken (
 ID BIGINT PRIMARY KEY AUTO_INCREMENT,
 UserID BIGINT NOT NULL,
 Token BINARY(60) NOT NULL,
 CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
 UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 IsProcessed BOOLEAN DEFAULT FALSE
);

CREATE INDEX user_id_idx ON ActivationToken (UserID);


