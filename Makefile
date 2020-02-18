###############################################
#
# Makefile
#
###############################################

.DEFAULT_GOAL := build

.PHONY: test

VERSION := 0.1.0

ver:
	@sed -i '' 's/^const Version = "[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}"/const Version = "${VERSION}"/' rx.go

lint:
	$(shell go env GOPATH)/bin/golint ./src/...

format:
	go fmt ./src/...

vet:
	go vet ./src/...

# PROFILE = -blockprofile
# PROFILE = -cpuprofile
PROFILE = -memprofile
profile:
	-rm -f rx.prof rx.test
	-go test ${PROFILE}=rx.prof ./src/...
	go tool pprof rx.prof

ppcpu:
	go tool pprof http://localhost/debug/pprof/profile

ppmem:
	go tool pprof http://localhost/debug/pprof/heap

build:
	go build -v ./...

clean:
	go clean ...

demo: build
	go run cmd/demo.go

race:
	go build -race ./...

test: build
	go test -v -count=1 ./...

bench: build
	go test -bench=. -v ./...

github:
	open "https://github.com/mlavergn/gorx"

release:
	zip -r gorx.zip LICENSE README.md Makefile cmd *.go go.mod
	hub release create -m "${VERSION} - GoRx" -a gorx.zip -t master "v${VERSION}"
	open "https://github.com/mlavergn/gorx/releases"
