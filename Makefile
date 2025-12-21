.PHONY: docker-build-mac
docker-build-mac:
	docker compose --profile build up
	docker build \
	--build-arg BUILD_GOOS=darwin \
	-t coutito/basement-mac:$(shell git rev-parse --short HEAD) \
	--target prod-build .

.PHONY: docker-build-aws
docker-build-aws:
	docker compose --profile build up
	docker build \
	--build-arg BUILD_GOOS=linux \
	-t coutito/basement-aws:$(shell git rev-parse --short HEAD) \
	--target prod-build .

.PHONY: docker-push-image
docker-push-image: docker-build-aws
	docker push coutito/basement-aws:$(shell git rev-parse --short HEAD)

.PHONY: docker-dev
docker-dev:
	docker compose --profile hot up -d
	docker logs -f basement-container
	
