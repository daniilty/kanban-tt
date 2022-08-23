package core

// Code - fail code.
type Code string

const (
	// все ок
	CodeOK = "OK"
	// внутренняя ошибка базы
	CodeDBFail = "DB_FUCKUP"
	// недостаточный доступ для операции
	CodeNotPermitted = "NOT_PERMITTED"
	// статус с таким айди не существует
	CodeNoStatus = "NO_SUCH_STATUS"
)
