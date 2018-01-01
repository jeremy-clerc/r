.PHONY: default test bindir all bundle install uninstall clean release

TAG := $(shell git describe --exact-match --tags HEAD 2>/dev/null)
SHA := $(shell git rev-parse --short HEAD)
DATE := $(shell date --iso-8601=seconds)
VERSION := $(TAG)
ifeq ($(strip $(VERSION)),)
	VERSION := $(SHA)
endif
PACKAGE := $(shell basename $(shell pwd))
BINARY_DIR := $(PACKAGE)-$(VERSION)
BINARIES := $(shell find cmd/ -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
.PHONY: $(BINARIES)

default:
	$(error no default action, but you can test, build any of $(BINARIES) or all, bundle, install, clean, release)

test:
	go test -v ./...

bindir:
	mkdir -p $(BINARY_DIR)/bin

$(BINARIES): bindir
	go build -ldflags "-X 'main.version=tag: $(TAG) commit: $(SHA) built: $(DATE)'" -o $(BINARY_DIR)/bin/$@ ./cmd/$@

all: $(BINARIES)

bundle: all
	tar --owner=root --group=root -cf $(BINARY_DIR).tar $(BINARY_DIR)/

install:
	sudo mkdir -p /usr/local/$(PACKAGE)
	sudo tar -xf $(BINARY_DIR).tar -C /usr/local/$(PACKAGE)
	if [ -e /usr/local/$(PACKAGE)/VERSION ]; then \
		sudo stow --dir /usr/local/$(PACKAGE) -D $(PACKAGE)-$(shell cat /usr/local/$(PACKAGE)/VERSION); \
	fi
	sudo stow --dir /usr/local/$(PACKAGE) -S $(PACKAGE)-$(VERSION)
	echo $(VERSION) | sudo tee /usr/local/$(PACKAGE)/VERSION >/dev/null

uninstall:
	sudo stow --dir /usr/local/$(PACKAGE) -D $(PACKAGE)-$(shell cat /usr/local/$(PACKAGE)/VERSION)
	sudo rm -r /usr/local/$(PACKAGE)

clean:
	rm -rf ./$(BINARY_DIR)/ $(BINARY_DIR).tar

release: test bundle install clean
	$(info  Released version $(VERSION))
