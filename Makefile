run: tidy
	go run example/main.go

test: tidy
	go test -v -cover

tidy:
	go mod tidy
