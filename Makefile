#TODO add .exe on windows boxen

BINDATA := $(GOPATH)/bin/go-bindata
BINDATA_DIRS = web/assets web/assets/css web/assets/js web/assets/fonts
BINDATA_FLAGS = -o=web/bindata.go -pkg=web -prefix web/assets $(BINDATA_DIRS)

.PHONY: start run deps bindata bindata-final compile release build

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

build: bindata-final compile

compile:
	go build github.com/cfstras/cfmedias/cmd/cfmedias

bindata-final:
	$(BINDATA) -debug=false -nocompress=false $(BINDATA_FLAGS)

bindata:
	$(BINDATA) -debug=true $(BINDATA_FLAGS)

run: build start

start:
	./cfmedias

clean:
	rm cfmedias

deps:
	@echo please install portaudio1.9-dev and libtagc0-dev with your package manager

	go get \
		code.google.com/p/go.tools/cmd/goimports \
		github.com/go-contrib/uuid \
		code.google.com/p/portaudio-go/portaudio \
		github.com/mattn/go-sqlite3 \
		github.com/jinzhu/gorm \
		github.com/peterh/liner \
		code.google.com/p/go.crypto/pbkdf2 \
		github.com/go-martini/martini \
		github.com/martini-contrib/render \
		github.com/cfstras/go-taglib
	go get github.com/jteeuwen/go-bindata/...
