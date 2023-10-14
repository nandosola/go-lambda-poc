# Go parameters
GOCMD=go
ZIPCMD=/usr/bin/zip
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=main
DIST_BINARY_NAME=bootstrap
DIST_ZIP=lambda-get
DIST_DIR=../../dist
TESTDATA_DIR=../../testdata
DOCKER_COMPOSE_YML=../../docker-compose.yaml

.PHONY: clean build dist run zip

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(DIST_BINARY_NAME)

build:
	$(GOCMD) mod tidy
	$(GOBUILD) -o $(BINARY_NAME) -v

dist: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath $(GOBUILD) -mod=readonly -ldflags='-s -w' -o $(DIST_BINARY_NAME) -v

zip: dist
	mv $(DIST_BINARY_NAME) $(DIST_DIR)
	cd $(DIST_DIR) && $(ZIPCMD) -j $(DIST_ZIP) $(DIST_BINARY_NAME) && rm $(DIST_BINARY_NAME)

run: dist
	docker-compose --file $(DOCKER_COMPOSE_YML) up -d --remove-orphans dynamo
	aws dynamodb create-table --no-cli-pager --cli-input-json file://$(TESTDATA_DIR)/create-ddb-table.json --endpoint-url http://localhost:8000
	aws dynamodb batch-write-item --request-items file://$(TESTDATA_DIR)/fixtures.json --endpoint-url http://localhost:8000
	bash -c "trap 'docker-compose --file $(DOCKER_COMPOSE_YML) down' EXIT; sam local start-api --docker-network lambda-local"

