CREATE TABLE task (
	ID INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
	STATUS VARCHAR(60) NOT NULL,
	TITLE VARCHAR(60) NOT NULL,
	ASIGNEE_ID VARCHAR(60) NOT NULL,
	CHANNEL_ID VARCHAR(60) NOT NULL
);