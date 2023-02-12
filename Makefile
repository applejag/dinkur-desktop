
.PHONY: build
build:
	wails build -tags='fts5' -ldflags='-s -w'

.PHONY: dev
dev:
	wails dev -tags='fts5'
