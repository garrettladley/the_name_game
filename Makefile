.PHONY: install gen-css watch-css install-templ gen-templ watch-templ build build-prod run

NODE_BIN := ./node_modules/.bin

install: install-templ gen-templ
	@go get ./...
	@go mod tidy
	@go mod download
	@mkdir -p htmx
	@wget -q -O htmx/htmx.min.js.gz https://unpkg.com/htmx.org@1.9.12/dist/htmx.min.js.gz
	@gunzip -f htmx/htmx.min.js.gz
	@npm install -D daisyui@latest
	@npm install -D tailwindcss

gen-css:
	@$(NODE_BIN)/tailwindcss build -i views/css/app.css -o public/styles.css --minify

watch-css:
	@$(NODE_BIN)/tailwindcss -i views/css/app.css -o public/styles.css --minify --watch 

install-templ:
	@go install github.com/a-h/templ/cmd/templ@latest

gen-templ:
	@templ generate

watch-templ:
	@templ generate --watch --proxy=http://127.0.0.1:3000

build: gen-css gen-templ
	@go build -tags dev -o bin/the_name_game main.go

build-prod: gen-css gen-templ
	@go build -tags prod -o bin/the_name_game main.go

run: build
	@./bin/the_name_game