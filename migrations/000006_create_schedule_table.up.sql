CREATE TABLE IF NOT EXISTS Schedule (
    ID BIGINT PRIMARY KEY AUTO_INCREMENT,
    ReminderID BIGINT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    DispatchTime DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    IsProcessed TINYINT DEFAULT 0,
    Version INT DEFAULT 1,
    FOREIGN KEY (ReminderID) REFERENCES Reminder(ID) ON DELETE CASCADE
);

CREATE INDEX dispatch_time_idx ON Schedule(DispatchTime);




INSERT INTO Schedule
(ReminderID, DispatchTime)
VALUES
(
    1,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 5 SECOND))
),
(
    1,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 3 MINUTE))
),
(
    1,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 5 MINUTE))
),
(
    2,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 2 MINUTE))
),
(
    2,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 15 MINUTE))
),
(
    3,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 WEEK))
),
(
    3,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 30 SECOND))
),
(
    3,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 MINUTE))
),
(
    3,
    TIMESTAMP(DATE_ADD(NOW(), INTERVAL 3 MINUTE))
);