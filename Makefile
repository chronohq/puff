VERSION ?= dev
BINARY_NAME = puff
RELEASE_DIR = dist
RELEASE_BUILD_FLAGS=-v -ldflags "-w -s -X main.version=$(VERSION)" -trimpath

build:
	go vet ./...
	go build -o $(BINARY_NAME)

release-darwin-arm64: require-version
	@echo "Building for darwin-arm64"
	rm -rf $(RELEASE_DIR)
	go vet ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${RELEASE_BUILD_FLAGS} -o $(RELEASE_DIR)/$(BINARY_NAME)
	tar -czvf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(RELEASE_DIR) $(BINARY_NAME)

release-linux-amd64: require-version
	@echo "Building for linux-amd64"
	rm -rf $(RELEASE_DIR)
	go vet ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${RELEASE_BUILD_FLAGS} -o $(RELEASE_DIR)/$(BINARY_NAME)
	tar -czvf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(RELEASE_DIR) $(BINARY_NAME)

release-windows-amd64:
	@echo "Building for windows-amd64"
	rm -rf $(RELEASE_DIR)
	go vet ./...
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${RELEASE_BUILD_FLAGS} -o $(RELEASE_DIR)/$(BINARY_NAME).exe
	zip -j $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(RELEASE_DIR)/$(BINARY_NAME).exe

release: release-darwin-arm64 release-linux-amd64 release-windows-amd64

require-version:
	@if [ "$(VERSION)" = "dev" ]; then \
		echo "You must set the VERSION"; \
		exit 1; \
	fi

clean:
	git clean -fd
