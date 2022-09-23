all: build generate

build: 
	go mod tidy
	go build -o bin/generate main.go

generate:
	./bin/generate