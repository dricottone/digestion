SHARE_DIR?=/usr/local/share
INSTALL_DIR?=$(SHARE_DIR)/mail-filters
GO_FILES=$(shell find -name '*.go')

.PHONY: clean
clean:
	rm -rf digestion go.sum

digestion: $(GO_FILES)
	go get -u
	go build

.PHONY: build
build: digestion

.PHONY: install
install: digestion
	mkdir -m755 -p $(INSTALL_DIR)
	install -m755 digestion $(INSTALL_DIR)/digestion

