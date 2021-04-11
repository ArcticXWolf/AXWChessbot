VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwchessbot

LDFLAGS=-ldflags "-w -s -X main.engineVersion=${VERSION} -X main.buildDate=${BUILD} -X main.gitCommit=${COMMIT}"

build:
	echo "Building for linux and windows"
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY} ${LDFLAGS} go.janniklasrichter.de/axwchessbot

run:
	go run go.janniklasrichter.de/axwchessbot

clean:
	rm -rf build/

all: clean build

.PHONY: clean build run all