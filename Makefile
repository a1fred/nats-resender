BRANCH=$(shell git branch --show-current)
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

TARGET=nats-resender

all: $(TARGET)

.PHONY: test
test:
	golangci-lint run
	go test -v -coverprofile=.coverage ./...

$(TARGET): test
	CGO_ENABLED=0 go build -ldflags "-X main.revision=$(REV) -s -w" -o $(TARGET) .

.PHONY: clean
clean:
	rm -rf $(TARGET) dist .coverage

.PHONY: dev nodev
dev:
	docker-compose -f docker-compose-dev.yml up -d --remove-orphans
nodev:
	docker-compose -f docker-compose-dev.yml rm -fs
