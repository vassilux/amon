# Makefile for the vorimport project
# Helper file
 
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w
GOTEST=go test -v
# Project target
TARGET=./bin/amon 
# Build source
SRC_BUILD_LIST = main.go
#Source files
SRC_LIST=$(wildcard *.go)
#Test files
TEST_LIST = $(wildcard *_test.go)
# 
BUILD_LIST = $(foreach int, $(SRC_BUILD_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(SRC_LIST), $(int)_clean)
FMT_TEST = $(foreach int, $(SRC_LIST), $(int)_fmt)
 
# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(BUILD_LIST)
 
all: build
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
test:  
	go test -v
fmt: $(FMT_TEST)
git:
	git commit -a -m "$m"
	git push https://github.com/vassilux/amon.git
depends:
	go get -u github.com/cihub/seelog
	go get -u github.com/cloudfoundry/gosigar
	go get -u github.com/bmizerany/pat
	
 
$(BUILD_LIST): %_build: %_fmt
	$(GOBUILD) -o $(TARGET)
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $*
$(TEST_LIST): %_test:
	$(GOTEST) $*
$(FMT_TEST): %_fmt:
	$(GOFMT) $*
