package db

const (
	CREATE_DB string = "CREATE DB loggerb"
	CREATE_TABLE_SENDER string = "CREATE TABLE IF NOT EXISTS Sender (senderID INTEGER,senderAddr TEXT NOT NULL,senderPort INTEGER NOT NULL,PRIMARY KEY(senderID));"
	CREATE_TABLE_MESSAGE string = "CREATE TABLE IF NOT EXISTS Message (messageID INTEGER,messageEndpoint TEXT NOT NULL,messageKey TEXT NOT NULL,messageValue TEXT NOT NULL,PRIMARY KEY(messageID));"
	CREATE_TABLE_BUFFER string =  "CREATE TABLE IF NOT EXISTS Buffer (logID INTEGER,logState INTEGER NOT NULL,senderID INTEGER NOT NULL,messageID INTEGER NOT NULL,PRIMARY KEY(logID),FOREIGN KEY (senderID) REFERENCES Sender(senderID),FOREIGN KEY (messageID) REFERENCES Message(messageID));"
	CREATE_TRIGGER string = "CREATE TRIGGER inserter AFTER INSERT ON Message BEGIN INSERT INTO Buffer (logState, senderID, messageID) VALUES (?, ?, ?); END;"
	INSERT_SENDER string = "INSERT INTO Sender VALUES (senderAddr, senderPort) VALUES (?, ?);"
	INSERT_MESSAGE string = "INSERT INTO Message VALUES (messageEndpoint, messageKey, messageValue) VALUES (?, ?);"
)
	
