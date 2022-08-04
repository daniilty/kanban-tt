package core

// Code - fail code.
type Code string

const (
	// все ок
	CodeOK = "OK"
	// внутренняя ошибка базы
	CodeDBFail = "DB_FUCKUP"
)
