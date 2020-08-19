SHARE_DIR?=/usr/local/share
INSTALL_DIR?=$(SHARE_DIR)/mail-filters

clean:
	rm -rf digestion

build:
	go build

install: build
	mkdir -m755 -p $(INSTALL_DIR)
	install -m755 digestion $(INSTALL_DIR)/digestion

