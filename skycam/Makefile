BASEDIR ?= ${PWD}
WORKDIR ?= $(PWD)/.work
HESTIA ?= ${BASEDIR}/hestia


.PHONY: all
all: auto test build

.PHONY: verbump
verbump:
	@APP_ROOT=${BASEDIR} ${HESTIA}/bin/auto_version


.PHONY: build
build:
	go build -o $(WORKDIR)/jupiter cmd/jupiter/main.go
	GOOS=linux GOARCH=arm64 go build -o ${WORKDIR}/jupiter_arm64 cmd/jupiter/main.go

.PHONY: run
run: build
	./.work/jupiter

.PHONY: auto
auto:
	go generate -x ./...

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
	rm -rf "${WORKDIR}"
	rm -rf dist
