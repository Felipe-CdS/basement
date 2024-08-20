#templ starts watching .templ files changes. Proxy argument points to the http port configured
#on the code.
live/templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx tailwindcss -i ./assets/css/tailwind.input.css -o ./assets/css/styles.css --minify --watch

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/main ./cmd/" --build.bin "./tmp/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

live: 
	make -j4 live/templ live/server live/sync_assets

build: 
	go build -o bin/main ./cmd
