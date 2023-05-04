package main

type Translations struct {
	DatabaseErrors struct {
		ConnectionFailed string
		MigrationFailed  string
		CreationFailed   string
		CloseFailed      string
	}
	CronJobErrors struct {
		FailedToStart string
	}
	HTTPErrors struct {
		FileTooLarge  string
		TorNotAllowed string
		BadRequest    string
	}
}

var en = Translations{
	// Errors related to the database.
	DatabaseErrors: struct {
		ConnectionFailed string
		MigrationFailed  string
		CreationFailed   string
		CloseFailed      string
	}{
		ConnectionFailed: "Failed to connect to the database.",
		MigrationFailed:  "Failed to migrate the database.",
		CreationFailed:   "Failed to create the database file.",
		CloseFailed:      "Failed to close the database file.",
	},

	// Errors related to the cron jobs (scheduled jobs).
	CronJobErrors: struct {
		FailedToStart string
	}{
		FailedToStart: "Failed to start the server auto-cleanup.",
	},

	// Errors related to the HTTP-server.
	HTTPErrors: struct {
		FileTooLarge  string
		TorNotAllowed string
		BadRequest    string
	}{
		FileTooLarge:  "The file you are trying to upload is too large.",
		TorNotAllowed: "You are not allowed to access the server from a Tor exit node.",
		BadRequest:    "Bad request.",
	},
}

var uk = Translations{
	// Errors related to the database.
	DatabaseErrors: struct {
		ConnectionFailed string
		MigrationFailed  string
		CreationFailed   string
		CloseFailed      string
	}{
		ConnectionFailed: "Не вдалося підключитися до бази даних.",
		MigrationFailed:  "Не вдалося мігрувати базу даних.",
		CreationFailed:   "Не вдалося створити файл бази даних.",
		CloseFailed:      "Не вдалося закрити файл бази даних.",
	},

	// Errors related to the cron jobs (scheduled jobs).
	CronJobErrors: struct {
		FailedToStart string
	}{
		FailedToStart: "Не вдалося запустити автоочистку сервера.",
	},

	// Errors related to the HTTP-server.
	HTTPErrors: struct {
		FileTooLarge  string
		TorNotAllowed string
		BadRequest    string
	}{
		FileTooLarge:  "Файл, який ви намагаєтесь завантажити, занадто великий.",
		TorNotAllowed: "Ви не можете отримати доступ до сервера з виходу Tor.",
		BadRequest:    "Поганий запит.",
	},
}

func translation(language string) Translations {
	switch language {
	case "uk":
		return uk
	default:
		return en
	}
}
