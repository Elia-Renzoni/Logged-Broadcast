/**
* SQL queries 
*/

CREATE DATABASE logger;

CREATE TABLE IF NOT EXISTS Buffer (
	logID INTEGER,
	logState INTEGER NOT NULL,
	senderID INTEGER NOT NULL,
	messageID INTEGER NOT NULL,
	CHECK (logState == 0 OR logState == 1),
	PRIMARY KEY(logID),
	CONSTRAINTS FOREIGN KEY Sender(senderID),
	CONSTRAINTS FOREIGN KEY Message(messageID)
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

/* select N pending message to reconstruct the 
 * pending buffer */
/** 
* this query will be triggered in the recovery
* session.
*/
/*
 * TODO -> end the following select query
*/
SELECT messageID, messageEndpoint, messageKey, messageValue
FROM Buffer
WHERE Buffer.messageID IN (
	SELECT messageID
	FROM Message
	GROUP BY messageKey
);
