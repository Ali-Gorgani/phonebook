# Makefile

.PHONY: all build clean

# Define the output directories for binaries
BIN_DIR := ./bin
PB_SERVER := $(BIN_DIR)/pbserver
PB_CLIENT := $(BIN_DIR)/pbclient

# Default target
all: build

# Build the pbserver and pbclient binaries
build: $(PB_SERVER) $(PB_CLIENT)

# Compile pbserver binary
$(PB_SERVER): cmd/pbserver/main.go
	mkdir -p $(BIN_DIR)
	go build -o $(PB_SERVER) ./cmd/pbserver

# Compile pbclient binary
$(PB_CLIENT): cmd/pbclient/main.go
	mkdir -p $(BIN_DIR)
	go build -o $(PB_CLIENT) ./cmd/pbclient

# Clean up binaries
clean:
	rm -rf $(BIN_DIR)
