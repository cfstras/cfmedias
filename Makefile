#TODO add .exe on windows boxen

BINDATA := $(GOPATH)/bin/go-bindata
BINDATA_DIRS = web/assets web/assets/css web/assets/js web/assets/fonts
BINDATA_FLAGS = -o=web/bindata.go -pkg=web -prefix web/assets $(BINDATA_DIRS)

.PHONY: all build build-debug compile bindata-final bindata-debug run start clean fix bindata-dep

FOLDERS := $(shell find * -type d)

all: build-debug

build: bindata-final compile
build-debug: bindata-debug compile

compile:
	@echo -------------------------------------------------------------------
	@echo if you encounter include errors please install portaudio1.9-dev and
	@echo libtagc0-dev with your package manager
	@echo -------------------------------------------------------------------
	go get -d
	go build -v

bindata-final: bindata-dep
	$(BINDATA) -debug=false -nocompress=false $(BINDATA_FLAGS)

bindata-debug: bindata-dep
	$(BINDATA) -debug=true $(BINDATA_FLAGS)

run: build-debug start

start:
	./cfmedias

clean:
	rm cfmedias

fix:
	goimports -l -w $(FOLDERS)
	for f in $(find -type f -name "*.go"); do go fix $i; done

bindata-dep:
	go get github.com/jteeuwen/go-bindata/...
