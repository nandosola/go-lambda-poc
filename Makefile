GOTIDY=go mod tidy -go=1.19 -compat=1.19
GOBUILD=$(GOTIDY) && go build -v
GOLIST=go list -m -json all
GOTEST=go test -v

.PHONY: build test list-deps deploy

build:
	cd service && $(GOBUILD)
	cd transport && $(GOBUILD)
	cd functions/get-birthday &&  $(GOBUILD)
	cd functions/put-birthday && $(GOBUILD)

list-deps: # useful for IDEs
	cd service && $(GOLIST)
	cd transport && $(GOLIST)
	cd functions/get-birthday &&  $(GOLIST)
	cd functions/put-birthday && $(GOLIST)

test:
	cd service && $(GOTEST)
	cd transport && $(GOTEST)

deploy:
	cd functions/get-birthday &&  $(MAKE) zip
	cd functions/put-birthday && $(MAKE) zip
	cd deploy && terraform plan -out myplan && terraform apply "myplan"

