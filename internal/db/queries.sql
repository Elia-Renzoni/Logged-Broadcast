/**
* SQL queries 
*/

CREATE DATABASE logger;

CREATE TABLE IF NOT EXISTS Buffer (
	logID INTEGER,
	logState INTEGER NOT NULL,
	senderID INTEGER NOT NULL,
	messageID INTEGER NOT NULL,
	PRIMARY KEY(logID),
	FOREIGN KEY (senderID) REFERENCES Sender(senderID),
	FOREIGN KEY (messageID) REFERENCES Message(messageID)
);

CREATE TABLE IF NOT EXISTS Sender (
	senderID INTEGER,
	senderAddr TEXT NOT NULL,
	senderPort INTEGER NOT NULL,
	PRIMARY KEY(senderID)
);

CREATE TABLE IF NOT EXISTS Message (
	messageID INTEGER,
	messageEndpoint TEXT NOT NULL,
	messageKey TEXT NOT NULL,
	messageValue TEXT NOT NULL,
	PRIMARY KEY(messageID)
);
