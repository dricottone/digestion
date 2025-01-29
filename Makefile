INSTALL_DIR?=/usr/local/bin
SRC=$(shell find -name '*.go')

clean:
	rm -f digestion go.mod go.sum

go.mod:
	go mod init git.sr.ht/~dricottone/digestion
	go get git.sr.ht/~dricottone/parcels
	go get git.sr.ht/~dricottone/textwrap

digestion: $(SRC)
	go get -u
	go build

build: go.mod digestion

install: digestion
	install --target-directory=$(INSTALL_DIR) --mode=755 digestion

uninstall:
	cd $(INSTALL_DIR) && rm -f digestion

.PHONY: clean build install uninstall
