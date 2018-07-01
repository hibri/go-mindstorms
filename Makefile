# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=legobot
BINARY_UNIX=$(BINARY_NAME)_linux_arm
	
default: test build

test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)



# Cross compilation
build:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_UNIX) -v