CREATE TABLE IF NOT EXISTS Reminder (
    ID BIGINT PRIMARY KEY AUTO_INCREMENT,
    Title VARCHAR(255) NOT NULL,
    Content TEXT,
    UserID BIGINT,
    DueDate DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    Version INT,
    FOREIGN KEY (UserID) REFERENCES User(ID) ON DELETE CASCADE
);  


INSERT INTO Reminder
(Title, Content, UserID, DueDate, Version)
VALUES
    (
        "Elizabeth, you are loved and worthy",
        "and i love you very much",
        1,
        TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 WEEK)),
        1
    ),
    (
        "I believe in you",
        "and I love you unconditionally",
        1,
        TIMESTAMP(DATE_ADD(NOW(), INTERVAL 2 WEEK)),
         1
    ),
    (
        "I care about you",
        "you enrich my life",
        1,
         TIMESTAMP(DATE_ADD(NOW(), INTERVAL 3 WEEK)),
        1
    );

