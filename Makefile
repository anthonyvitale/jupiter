test:
	go test ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

gen:
	go generate ./...

clean:
	go clean
	rm coverage.out
