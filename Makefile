GIT_SHA := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags --abbrev=0)
BUILD_TIME := $(shell date -Iseconds -u)

.PHONY: build
build:
	docker buildx build \
	--push \
	--no-cache \
	-t coutito/basement:latest \
	--platform linux/amd64,linux/arm64 \
	--build-arg GIT_TAG=$(GIT_TAG) \
	--build-arg GIT_SHA=$(GIT_SHA) \
	--build-arg BUILD_TIME=$(BUILD_TIME) \
	-t coutito/basement:$(GIT_SHA) .

.PHONY: run
run:
	docker run -d -p 3000:3000 \
		-v .env:/app/.env \
		coutito/basement:$(GIT_SHA)

.PHONY: dev
dev:
	GIT_TAG=$(GIT_TAG) GIT_SHA=$(GIT_SHA) BUILD_TIME=$(BUILD_TIME) docker compose --profile hot up -d 
	docker logs -f basement-container
