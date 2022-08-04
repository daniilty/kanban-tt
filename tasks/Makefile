build:
	go build github.com/daniilty/kanban-tt/tasks/cmd/server
build_docker:
	docker build -t kanban-tt/tasks-api:latest -f docker/Dockerfile .
gen_swagger:
	swagger generate spec -o ./api/openapi/swagger.json -m
