all: build
GOFILES=`go list ./... | grep -v vendor`

test:
	go test $(GOFILES)

build: test
	go build

escape:
	go build --gcflags="-m -m" .