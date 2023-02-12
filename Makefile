PREFIX=$$HOME/.local
# Alternative prefixes:
# PREFIX=/usr
BIN_PREFIX=${PREFIX}/bin
ICON_PREFIX=${PREFIX}/share/icons
APPS_PREFIX=${PREFIX}/share/applications

.PHONY: build
build: build/bin/dinkur-desktop

build/bin/dinkur-desktop: $(shell git ls-files '*.go') Makefile
	wails build -tags='fts5' -ldflags='-s -w'

.PHONY: install
install:
	cp build/bin/dinkur-desktop ${BIN_PREFIX}/dinkur-desktop
	mkdir -pv ${ICON_PREFIX}/hicolor/64x64/apps
	mkdir -pv ${ICON_PREFIX}/hicolor/48x48/apps
	mkdir -pv ${ICON_PREFIX}/hicolor/scalable/apps
	cp icons/dinkur-small-64.png ${ICON_PREFIX}/hicolor/64x64/apps/dinkur-small.png
	cp icons/dinkur-small-48.png ${ICON_PREFIX}/hicolor/48x48/apps/dinkur-small.png
	cp icons/dinkur-small.svg ${ICON_PREFIX}/hicolor/scalable/apps/dinkur-small.svg
	desktop-file-install build/linux/dinkur-desktop.desktop --dir=${APPS_PREFIX}

.PHONY: dev
dev:
	wails dev -tags='fts5'
