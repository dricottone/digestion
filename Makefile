INSTALL_DIR?=/usr/local/bin
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
	install -m755 digestion $(INSTALL_DIR)/digestion

