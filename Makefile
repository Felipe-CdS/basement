#templ starts watching .templ files changes. Proxy argument points to the http port configured
#on the code.
live/templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false


# run air to detect any go file changes to re-build and re-run the server.
live/server:
	air \
	--build.cmd "go build -o tmp/main ./cmd/" --build.bin "./tmp/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx tailwindcss -i ./assets/css/tailwind.input.css -o ./assets/css/tailwind_styles.css --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	air \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

live: 
	make -j3 live/templ live/server live/sync_assets

build_deprecated:
	go run github.com/a-h/templ/cmd/templ@v0.3.819 generate
	tailwindcss -i ./assets/css/tailwind.input.css -o ./assets/css/tailwind_styles.css --minify
	go build -o bin/main ./cmd

.PHONY: docker-build-mac
docker-build-mac:
	docker build \
	--build-arg BUILD_GOOS=darwin \
	-t basement:$(shell git rev-parse --short HEAD) \
	--target prod-build .

.PHONY: docker-build-aws
docker-build-aws:
	docker build \
	--build-arg BUILD_GOOS=linux \
	-t basement:aws \
	-t basement:$(shell git rev-parse --short HEAD) \
	--target prod-build .

.PHONY: docker-push-image
docker-push-image:
	docker image tag basement:aws coutito/basement:aws
	docker push coutito/basement:aws

.PHONY: docker-dev
docker-dev:
	docker build -t basement:dev --target dev-hot .
	docker stop basement > /dev/null 2>&1 || true
	docker rm basement   > /dev/null 2>&1 || true
	docker run -d -p 5432:8000 --name basement -v ${PWD}:/app basement:dev
	
