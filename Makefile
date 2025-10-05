BINARY := pipause
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin

.PHONY: all build install uninstall test

all: build

build:
	go build -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) $(BINDIR)/$(BINARY)

uninstall:
	rm -f $(BINDIR)/$(BINARY)

test:
	@echo "Running tests (none yet — placeholder)"
	@go test ./...

