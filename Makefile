run: build
	@./bin/the_name_game

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D daisyui@latest
	@npm install -D tailwindcss

css:
	@tailwindcss -i views/css/app.css -o public/styles.css --watch 

templ:
	@templ generate --watch --proxy=http://127.0.0.1:3000

build:
	@templ generate views
	@go build -tags dev -o bin/the_name_game main.go 
