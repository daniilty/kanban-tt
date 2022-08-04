package server

import "github.com/daniilty/kanban-tt/tasks/internal/core"

const (
	// нет sub в заголовке authorization: Bearer xxx
	codeUnauthorizedNoSub core.Code = "NO_SUBJECT"

	// айди должен быть положительным целочисленным
	codeIDPositive core.Code = "ID_MUST_BE_POSITIVE"
	// нет айди
	codeNoID core.Code = "NO_ID"
	// айди должен быть целочисленным
	codeInvalidIDType core.Code = "INVALID_ID_TYPE"
	// пришло пустое тело запроса
	codeEmptyBody core.Code = "EMPTY_BODY"
	// необрабатываемые данные
	codeInvalidBodyStructure core.Code = "INVALID_BODY_STRUCTURE"
)
