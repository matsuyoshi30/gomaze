BIN := gomaze
export GO111MODULE=on

.PHONY: all
all: clean build

.PHONY: build
build:
	go build -o build/$(BIN)

.PHONY: clean
clean:
	rm -rf build
	go clean

.PHONY: test
test: build
	go test -v
