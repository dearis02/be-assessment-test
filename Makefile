# Usage: make migration:new name=init_users
migration\:new:
	migrate create -ext sql -dir ./migrations ${name}

# Usage: make migration:up database=postgres://user:password@127.0.0.1:5432/db_name?sslmode=disable
migration\:up:
	migrate -database ${database} -path ./migrations up

# Usage: make migration:down database=postgres://user:password@127.0.0.1:5432/db_name?sslmode=disable
migration\:down:
	migrate -database ${database} -path ./migrations down

# Usage: make docker:compose-up env=.env
docker\:compose-up:
	docker compose --env-file "${env}" up -d