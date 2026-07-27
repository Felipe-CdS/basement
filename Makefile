.PHONY: build
build:
	docker buildx build \
	--push \
	--no-cache \
	--platform linux/amd64,linux/arm64 \
	-t coutito/basement:$(shell git rev-parse --short HEAD) .

.PHONY: run
run:
	docker run -d -p 3000:3000 \
		-v .env:/app/.env \
		coutito/basement:$(shell git rev-parse --short HEAD)

.PHONY: dev
dev:
	docker compose --profile hot up -d
