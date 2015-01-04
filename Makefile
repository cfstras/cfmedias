#TODO add .exe on windows boxen

BINDATA := go-bindata
BINDATA_DIRS = $(shell find web/assets -type d)
BINDATA_FLAGS = -o=web/bindata.go -pkg=web -prefix web/assets $(BINDATA_DIRS)

.PHONY: all build build-debug compile bindata-final bindata-debug run start clean fix bindata-dep

FOLDERS := $(shell find * -type d)

all: build-debug

build: bindata-final compile
build-debug: bindata-debug compile

compile:
	@echo -------------------------------------------------------------------
	@echo if you encounter include errors please install portaudio1.9-dev,
	@echo libtagc0-dev and libgpod-dev with your package manager
	@echo -------------------------------------------------------------------
	go get -d
	go build -v

bindata-final: bindata grunt
	$(BINDATA) -debug=false -nocompress=false $(BINDATA_FLAGS)

bindata-debug: bindata grunt
	$(BINDATA) -debug=true $(BINDATA_FLAGS)

run: build-debug start

test:
	go get -d ./...
	go test ./...

start:
	./cfmedias

clean:
	rm -f cfmedias
	rm -rf web/node_modules
	rm -rf web/bower_components
	rm -rf web/assets/vendor
	rm -rf web/bindata.go
	cd web/list-view && git clean -dfx

fix: goimports
	goimports -l -w $(FOLDERS)
	for f in $$(find . -type f -name "*.go"); do \
		go fix "$$f"; \
		go tool vet -composites=false "$$f"; \
	done

.PHONY: grunt
grunt: web/assets/vendor

web/assets/vendor: web/bower_components web/Gruntfile.js web/list-view/dist/list-view.js
	cd web && node node_modules/grunt-cli/bin/grunt
	touch web/assets/vendor

web/list-view/dist/list-view.js:
	cd web/list-view && npm install && npm run build-all

web/bower_components: web/node_modules web/bower.json
	cd web && node node_modules/bower/bin/bower install
	touch web/bower_components

web/node_modules: web/package.json
	cd web && npm install --quiet
	touch web/node_modules

goimports:
	go get -v code.google.com/p/go.tools/cmd/goimports

bindata:
	go get -v github.com/jteeuwen/go-bindata/...
