PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

FUZZTIME ?= 10s

.PHONY: all build test test-race fuzz vet bench check-all install clean

all: build

build:
	cd go-notify && go build -o go-notify .
	cd go-notify-manager && go build -o go-notify-manager .

test:
	cd go-notify && go test ./...
	cd go-notify-manager && go test ./...

test-race:
	cd go-notify && go test -race ./internal/dbus ./internal/repository ./internal/service ./internal/viewmodel
	cd go-notify-manager && go test -race ./internal/extension ./internal/repository ./internal/services ./internal/viewmodel

vet:
	cd go-notify && go vet ./...
	cd go-notify-manager && go vet ./...

fuzz:
	cd go-notify && go test -fuzz=FuzzExtractURLs -fuzztime=$(FUZZTIME) ./internal/service
	cd go-notify && go test -fuzz=FuzzCleanExecLine -fuzztime=$(FUZZTIME) ./internal/service
	cd go-notify-manager && go test -fuzz=FuzzLinkify -fuzztime=$(FUZZTIME) ./internal/ui
	cd go-notify-manager && go test -fuzz=FuzzResolveAppIconName -fuzztime=$(FUZZTIME) ./internal/ui
	cd go-notify-manager && go test -fuzz=FuzzDateFormatToHMS -fuzztime=$(FUZZTIME) ./internal/extension

bench:
	cd go-notify && go test -bench=. -benchmem ./...
	cd go-notify-manager && go test -bench=. -benchmem ./...

check-all: vet test test-race bench

install: build
	mkdir -p $(BINDIR)
	cp -f go-notify/go-notify $(BINDIR)/go-notify
	cp -f go-notify-manager/go-notify-manager $(BINDIR)/go-notify-manager

clean:
	rm -f go-notify/go-notify go-notify-manager/go-notify-manager
