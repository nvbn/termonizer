build:
	go build -o bin/termonizer cmd/termonizer/main.go

test:
	go test -v ./...

run:
	go run cmd/termonizer/main.go
