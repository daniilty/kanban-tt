build:
	go build github.com/daniilty/kanban-tt/auth/cmd/server
build_docker:
	docker build -t kanban-tt/auth-api:latest -f docker/Dockerfile .
gen_swagger:
	swagger generate spec -o ./api/openapi/swagger.json -m

