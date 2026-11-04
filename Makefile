VERSION=v0.0.1-dev
CONFIG=
PORT=
ROOT=
OUTPUT=
URL=
ARGS=$(if $(CONFIG),-config=$(CONFIG)) $(if $(ROOT),-root=$(ROOT)) $(if $(URL),-url=$(URL))
SHELL := /bin/bash
.SHELLFLAGS := -e -o pipefail -c

build: clean
	go build -ldflags="-X main.Version=$(VERSION)" -o build/qwc .

pull:
	git pull origin main

package:
	mkdir -p dist
	zip -j dist/qwc-$(VERSION).zip build/qwc example.qwcconfig.json README.md LICENSE

clean:
	rm -rf build
	rm -rf dist

clear:
	clear

release: build package

move:
	mv build/qwc .

run-latest: pull run

run: build move clean clear
	./qwc $(ARGS) $(if $(PORT),-port=$(PORT))

site: build move clean
	./qwc $(ARGS) -build $(if $(OUTPUT),-output=$(OUTPUT))
