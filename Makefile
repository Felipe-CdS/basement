-include Makefile.vars

.PHONY: docker-build-mac
docker-build-mac:
	docker compose --profile build up
	docker rm tailwind-minify-container
	docker build \
	--build-arg BUILD_GOOS=darwin \
	-t coutito/basement-mac:$(shell git rev-parse --short HEAD) \
	--target prod-final-bin .

.PHONY: docker-build-aws
docker-build-aws:
	docker compose --profile build up
	docker rm tailwind-minify-container
	docker build \
	--build-arg BUILD_GOOS=linux \
	-t coutito/basement-aws:$(shell git rev-parse --short HEAD) \
	--target prod-final-bin .

.PHONY: docker-push-image
docker-push-image: docker-build-aws
	docker push coutito/basement-aws:$(shell git rev-parse --short HEAD)

.PHONY: docker-dev
docker-dev:
	docker compose --profile hot up -d

.PHONY: run-act
run-act:
	@if [ -z "$$AWS_PROFILE" ]; then \
	  echo "AWS_PROFILE is undefined"; exit 1; \
	fi

	act \
	  -s REMOTE_HOST="$(REMOTE_HOST)" \
	  -s REMOTE_USER="$(REMOTE_USER)" \
	  -s SSH_PRIVATE_KEY="$(PEM_KEY)" \
	  -P ubuntu-latest=catthehacker/ubuntu:act-latest
