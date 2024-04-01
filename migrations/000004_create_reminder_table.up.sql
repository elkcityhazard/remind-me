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
        "Go To Dentist",
        "You Have A Dentist Appointment With Dr. Doom",
        1,
        TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 WEEK)),
        1
    ),
    (
        "Go To Doctor",
        "You have to get your foot removed",
        1,
        TIMESTAMP(DATE_ADD(NOW(), INTERVAL 2 WEEK)),
         1
    ),
    (
        "Finish Book Report",
        "Must finish book report on 1984",
        1,
         TIMESTAMP(DATE_ADD(NOW(), INTERVAL 3 WEEK)),
        1
    );

