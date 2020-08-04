# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
AGENT_BINARY_NAME=icinga2-agent
CHECK_BINARY_NAME=check_g3000
OUT_DIR=bin

all: test build-agent build-check

build-agent:
	$(info $(shell mkdir -p $(OUT_DIR)))
	CC=/usr/local/musl/bin/musl-gcc $(GOBUILD) -o ./$(OUT_DIR)/$(AGENT_BINARY_NAME) --ldflags '-linkmode external -extldflags "-static"' ./agent/agent.go
	sha256sum ./$(OUT_DIR)/$(AGENT_BINARY_NAME) > ./$(OUT_DIR)/$(AGENT_BINARY_NAME).sha256

build-check: build-check.linux.amd64 build-check.linux.arm64 build-check.linux.arm5 build-check.linux.arm6 build-check.linux.arm7 build-check.windows.amd64 build-check.darwin.amd64

build-check.linux.amd64:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=linux GOARCH=amd64 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.amd64" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.amd64 > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.amd64.sha256

build-check.linux.arm64:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=linux GOARCH=arm64 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm64" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm64 > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm64.sha256

build-check.linux.arm5:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=linux GOARCH=arm GOARM=5 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm5" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm5 > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm5.sha256

build-check.linux.arm6:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=linux GOARCH=arm GOARM=6 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm6" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm6 > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm6.sha256

build-check.linux.arm7:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=linux GOARCH=arm GOARM=7 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm7" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm7 > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).linux.arm7.sha256

build-check.windows.amd64:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=windows GOARCH=amd64 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).windows.amd64.exe" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).windows.amd64.exe > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).windows.amd64.exe.sha256

build-check.darwin.amd64:
	$(info $(shell mkdir -p $(OUT_DIR)))
	GOOS=darwin GOARCH=amd64 go build -a -v -o "./$(OUT_DIR)/$(CHECK_BINARY_NAME).darwin.amd64" ./check/
	sha256sum ./$(OUT_DIR)/$(CHECK_BINARY_NAME).darwin.amd64 > ./$(OUT_DIR)/$(CHECK_BINARY_NAME).darwin.amd64.sha256

test:
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(OUT_DIR)
deps:
	$(GOGET) github.com/fatih/structs
	$(GOGET) github.com/mackerelio/go-osstat
	$(GOGET) github.com/mitchellh/mapstructure
	$(GOGET) github.com/urfave/cli/v2
