
TARGET_BIN = meta-generator
TARGET_ARCH = amd64
SOURCE_MAIN = main.go
LDFLAGS = -s -w

all: build

build: build-darwin build-linux build-windows

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_darwin-amd64 $(SOURCE_MAIN)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_linux-amd64 $(SOURCE_MAIN)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=$(TARGET_ARCH) go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET_BIN)_windows-amd64.exe $(SOURCE_MAIN)

start:
	go run $(SOURCE_MAIN) -views=./views

docker-build:
	docker build . -t robbymilo/meta-generator

run:
	docker run -p 3000:3000 --name=meta --rm robbymilo/meta-generator