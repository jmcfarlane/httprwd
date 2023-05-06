test: tidy
	go test -v -cover -race

tidy:
	go mod tidy
