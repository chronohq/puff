BINARY_NAME=puff

build:
	go vet ./...
	go build -o $(BINARY_NAME)

clean:
	git clean -fd