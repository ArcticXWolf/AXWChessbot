VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwchessbot

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.date=${BUILD} -X main.commit=${COMMIT}"

build:
	echo "Building for linux and windows"
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 ${LDFLAGS} go.janniklasrichter.de/axwchessbot
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe ${LDFLAGS} go.janniklasrichter.de/axwchessbot

run:
	go run go.janniklasrichter.de/axwchessbot

profile: clean
	mkdir build
	go test -count=10 -run=^$$ -bench "^(BenchmarkSearchFullEvaluation[3-5])$$" -benchmem -o build/test.bin -cpuprofile build/cpu.out -memprofile build/mem.out go.janniklasrichter.de/axwchessbot/search
	go tool pprof --svg build/test.bin build/mem.out > build/mem.svg
	go tool pprof --svg build/test.bin build/cpu.out > build/cpu.svg

clean:
	rm -rf build/

all: clean build

.PHONY: clean build run all