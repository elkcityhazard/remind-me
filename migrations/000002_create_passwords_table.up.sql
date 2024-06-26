
CREATE TABLE IF NOT EXISTS Password (
    ID BIGINT PRIMARY KEY AUTO_INCREMENT,
    Hash BLOB,
    UserID BIGINT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    Version INT,
    FOREIGN KEY (UserID) REFERENCES User(ID) ON DELETE CASCADE
);

