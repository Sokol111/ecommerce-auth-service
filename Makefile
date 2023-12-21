compile:
	mockery

docker-compose-up:
	docker compose up -d

docker-compose-down:
	docker compose down

build-docker-image:
	docker build -t sokol111/ecommerce-auth-service:latest .

test:
	go test ./... -v -cover
