VERSION ?= dev
BINARY_NAME = puff

build:
	go vet ./...
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)

clean:
	git clean -fd