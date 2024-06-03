install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod tidy
	@go mod download
	@mkdir -p htmx
	@wget -O htmx/htmx.min.js.gz https://unpkg.com/htmx.org@1.9.12/dist/htmx.min.js.gz
	@gunzip htmx/htmx.min.js.gz
	@npm install -D daisyui@latest
	@npm install -D tailwindcss

gen-css:
	@tailwindcss build -i views/css/app.css -o public/styles.css

watch-css:
	@tailwindcss -i views/css/app.css -o public/styles.css --watch 

gen-templ:
	@templ generate views

watch-templ:
	@templ generate --watch --proxy=http://127.0.0.1:3000

build: gen-css gen-templ
	@go build -tags dev -o bin/the_name_game main.go 

run: build
	@./bin/the_name_game