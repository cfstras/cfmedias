
#TODO add .exe on windows boxen

GOPATH := $(CURDIR)
export GOPATH

BINDATA := bin/go-bindata

FOLDERS := src/config \
	src/core \
	src/coreimpl \
	src/db \
	src/errrs \
	src/logger \
	src/main \
	src/util \
	src/web

all: build

release: bindata compile
build: bindata-debug compile

.PHONY: compile
compile:
	mkdir -p bin
	go build -o bin/cfmedias main

.PHONY: fix
fix:
	goimports -l -w $(FOLDERS)

BINDATA_DIRS = src/web/assets src/web/assets/css src/web/assets/js src/web/assets/fonts
BINDATA_FLAGS = -o=src/web/bindata.go -pkg=web -prefix src/web/assets $(BINDATA_DIRS)

.PHONY: bindata
bindata:
	$(BINDATA) -debug=false -nocompress=false $(BINDATA_FLAGS)

.PHONY: bindata-debug
bindata-debug:
	$(BINDATA) -debug=true $(BINDATA_FLAGS)


.PHONY: clean
clean:
	rm -rf bin

run: build start

start:
	bin/cfmedias

deps:
	go get github.com/mattn/go-sqlite3 \
	github.com/go-contrib/uuid \
	code.google.com/p/portaudio-go/portaudio \
	github.com/peterh/liner \
	code.google.com/p/go.crypto/pbkdf2 \
	github.com/jteeuwen/go-bindata/... \
	github.com/go-martini/martini \
	github.com/martini-contrib/render \
	github.com/jinzhu/gorm

devdeps: deps
	go get code.google.com/p/go.tools/cmd/goimports

	@echo please install portaudio1.9-dev and libtagc0-dev with your package manager
