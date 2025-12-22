.PHONY: docker-build-mac
docker-build-mac:
	docker compose --profile build up
	docker rm tailwind-minify-container
	docker rmi coutito/tailwindcss:v4.1.18
	docker build \
	--build-arg BUILD_GOOS=darwin \
	-t coutito/basement-mac:$(shell git rev-parse --short HEAD) \
	--target prod-final-bin .

.PHONY: docker-build-aws
docker-build-aws:
	docker compose --profile build up
	docker rm tailwind-minify-container
	docker rmi coutito/tailwindcss:v4.1.18
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
	act -s REMOTE_HOST=56.125.128.202 -s REMOTE_USER=rapi -s SSH_PRIVATE_KEY="$(shell cat ~/.ssh/gureum-key.pem)" -P ubuntu-latest=catthehacker/ubuntu:act-latest

