CREATE TABLE IF NOT EXISTS Users (
	ID INTEGER AUTO_INCREMENT,
	CreationTime DATETIME,
	LastUpdatedTime DATETIME,
	DeletionTime DATETIME,
	Status INTEGER,
	Username VARCHAR(128),
	Password CHAR(60),
	Email VARCHAR(254),
	CONSTRAINT pk_UsersID PRIMARY KEY (ID),
	CONSTRAINT uc_UsersUsername UNIQUE (Username)
);