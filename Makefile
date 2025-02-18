build:
	go build -o bin/termonizer cmd/termonizer/main.go

test:
	go test -v ./...

run:
	go run cmd/termonizer/main.go -debug debug.log

generate-test-db:
	go run cmd/generate-lorem-ipsum-db/main.go

run-test-db:
	go run cmd/termonizer/main.go -db test.db -debug debug.log

log:
	tail -f debug.log

sign-macos:
	cd bin && \
	codesign -s $$(security find-identity -v -p codesigning | grep 'Developer ID Application' | awk '{ print $$2 }') -o runtime -v termonizer && \
	zip termonizer.zip termonizer && \
	xcrun notarytool submit termonizer.zip --keychain-profile "notarization-profile" --wait
