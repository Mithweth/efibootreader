GO ?= go
REMOVE ?= rm
INSTALLBIN ?= install

ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

default: build

build: fmt
	@$(GO) mod download
	@$(GO) build -o efibootreader ./cmd

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

lint:
	@golangci-lint run

test: fmt vet lint
	@$(GO) test ./... -coverprofile=cover.out

coverage:
	@$(GO) tool cover -func=cover.out

clean:
	@$(REMOVE) -f efibootreader cover.out

install: 
	$(INSTALLBIN) -d $(PREFIX)/bin/
	$(INSTALLBIN) -m 755 efibootreader $(PREFIX)/bin/

uninstall:
	$(REMOVE) -f $(PREFIX)/bin/efibootreader
