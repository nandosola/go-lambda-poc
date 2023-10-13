# Go parameters
GOCMD=go
ZIPCMD=/usr/bin/zip
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=main
DIST_BINARY_NAME=bootstrap
DIST_ZIP=lambda-get
DIST_DIR=dist

.PHONY: clean build dist

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(DIST_BINARY_NAME)

dist: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath $(GOBUILD) -mod=readonly -ldflags='-s -w' -o $(DIST_BINARY_NAME) -v

zip: dist
	mv $(DIST_BINARY_NAME) $(DIST_DIR)
	cd $(DIST_DIR) && $(ZIPCMD) -j $(DIST_ZIP) $(DIST_BINARY_NAME) && rm $(DIST_BINARY_NAME)

run: dist
	sam local start-api

