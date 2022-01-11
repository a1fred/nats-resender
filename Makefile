BRANCH=$(shell git branch --show-current)
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)
TARGET=nats-resender

all: dev test $(TARGET)

$(TARGET): $(wildcard src/**/*)
	CGO_ENABLED=0 go build -ldflags "-X main.revision=$(REV) -s -w" -o $(TARGET) .

.PHONY: clean
clean:
	rm -rf $(TARGET) .coverage certs

.PHONY: test
test:
	golangci-lint run
	go test -v -coverprofile=.coverage ./...

certs:
	mkdir -p certs
	openssl genrsa -out certs/server.key 4096
	openssl req -new -x509 -sha256 -key certs/server.key -subj '/CN=server' -out certs/server.crt -days 3650 -addext 'subjectAltName = IP:127.0.0.1'

.PHONY: dev nodev
dev: certs
	docker-compose -f docker-compose-dev.yml up -d --remove-orphans
nodev:
	docker-compose -f docker-compose-dev.yml rm -fs
