.PHONY: generate css build run test check tidy clean

TEMPL := $(shell go env GOPATH)/bin/templ

generate:
	$(TEMPL) generate

css:
	pnpm css

build: generate css
	go build -o bin/gonshin ./cmd/gonshin

run: build
	./bin/gonshin

test: generate
	go test ./...

check: generate
	go vet ./...
	go test ./...
	go build ./...

tidy:
	go mod tidy

clean:
	rm -rf bin web/static/app.css
