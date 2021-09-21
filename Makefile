BRANCH=$(shell git branch --show-current)
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

BUILD_DIR=build
BIN=nats-resender
TARGET=$(BUILD_DIR)/$(BIN)

all: $(TARGET)

.PHONY: test
test:
	go test .

.PHONY: $(TARGET)
$(TARGET): test
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags "-X main.revision=$(REV) -s -w" -o $(TARGET) ./

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
