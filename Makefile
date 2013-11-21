
#TODO add .exe on windows boxen

GOPATH := $(CURDIR)

all: build

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/cfmedias -race main

.PHONY: clean
clean:
	rm -rf bin

run: build start

start:
	bin/cfmedias
