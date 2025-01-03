build:
	go build -o bin/termonizer cmd/termonizer/main.go

test:
	go test -v ./...

run:
	go run cmd/termonizer/main.go

generate-test-db:
	go run cmd/generate-lorem-ipsum-db/main.go

run-test-db:
	go run cmd/termonizer/main.go -db test.db
