PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

.PHONY: all build install clean

all: build

build:
	cd go-notify && go build -o go-notify .
	cd go-notify-manager && go build -o go-notify-manager .

install: build
	mkdir -p $(BINDIR)
	cp -f go-notify/go-notify $(BINDIR)/go-notify
	cp -f go-notify-manager/go-notify-manager $(BINDIR)/go-notify-manager

clean:
	rm -f go-notify/go-notify go-notify-manager/go-notify-manager
