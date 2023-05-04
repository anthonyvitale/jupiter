BASEDIR ?= ${PWD}
WORKDIR ?= $(PWD)/.work


.PHONY: all
all: auto test build

.PHONY: auto
auto:
	go generate -x ./...

.PHONY: build
build: auto
	go build -o $(WORKDIR)/jupiter cmd/cli/main.go

.PHONY: run
run: build
	./.work/jupiter

.PHONY: test
test: auto
	go test ./...

.PHONY: cov
cov: auto
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	find "${BASEDIR}" -name "*.auto.go" -print | xargs rm -f
	go clean
	rm -f "${BASEDIR}/coverage.out"
