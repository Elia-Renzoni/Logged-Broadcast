package db

var createTableSender string = 
`
CREATE TABLE IF NOT EXISTS Sender (
 senderID INTEGER,
 senderAddr TEXT NOT NULL,
 senderPort INTEGER NOT NULL,
 PRIMARY KEY(senderID)
);
`

var createTableMessage string = 
`
CREATE TABLE IF NOT EXISTS Message (
 messageID INTEGER,
 messageEndpoint TEXT NOT NULL,
 messageKey TEXT NOT NULL,
 messageValue TEXT NOT NULL,
 PRIMARY KEY(messageID)
);
`

var createTableBuffer string = 
`
CREATE TABLE IF NOT EXISTS Buffer (
 logID INTEGER,
 logState INTEGER NOT NULL,
 senderID INTEGER NOT NULL,
 messageID INTEGER NOT NULL,
 PRIMARY KEY(logID),
 FOREIGN KEY(senderID) REFERENCES Sender(senderID),
 FOREIGN KEY(messageID) REFERENCES Message(messageID)
);
`

var insertBufferStmt string = 
`
INSERT INTO Buffer (logState, senderID, messageID) VALUES (?, ?, ?);
`

var insertSenderStmt string = 
`
INSERT INTO Sender (senderAddr, senderPort) VALUES (?, ?);
`

var insertMessageStmt string = 
`
INSERT INTO Message (messageEndpoint, messageKey, messageValue)
VALUES (?, ?, ?);
`

var deleteMessageStmt string = 
`
DELETE FROM Message
WHERE messageKey = ?;
`

var deleteSenderStmt string = 
`
DELETE FROM Sender
WHERE senderID = ?;
`

var fetchMessageIDStmt string = 
`
SELECT Message.messageID
FROM Message
WHERE Message.messageKey = ?;
`

var fetchSenderIDStmt string =
`
SELECT Buffer.senderID
FROM Buffer
WHERE Buffer.messageID = ?;
`

var deleteEntriesFromBufferStmt string = 
`
DELETE FROM Buffer
WHERE messageID = ? AND senderID = ?;
`

var changeStatusToDelivered string = 
`
UPDATE Buffer
SET logState = 1
WHERE Buffer.messageID = ?;
`
