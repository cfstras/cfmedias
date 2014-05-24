
#TODO add .exe on windows boxen

GOPATH := $(CURDIR)
export GOPATH

BINDATA := bin/go-bindata

all: build

release: bindata compile
build: bindata-debug compile

.PHONY: compile
compile:
	mkdir -p bin
	go build -o bin/cfmedias main

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
	go get github.com/mattn/go-sqlite3
	go get github.com/go-contrib/uuid
	go get code.google.com/p/portaudio-go/portaudio
	go get github.com/coopernurse/gorp
	go get github.com/peterh/liner
	go get code.google.com/p/go.crypto/pbkdf2
	go get github.com/jteeuwen/go-bindata/...

	@echo please install portaudio1.9-dev and libtagc0-dev with your package manager
