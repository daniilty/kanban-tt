build_images:
	docker build -t kanban-tt/auth-api:latest -f ./auth/docker/Dockerfile ./auth
	docker build -t kanban-tt/tasks-api:latest -f ./tasks/docker/Dockerfile ./tasks
	docker build -t kanban-tt/users-api:latest -f ./users/docker/Dockerfile ./users
	docker build -t kanban-tt/email:latest -f ./email/docker/Dockerfile ./email
