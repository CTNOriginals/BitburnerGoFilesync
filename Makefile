# Makefile template: https://gist.github.com/grihabor/4a750b9d82c9aa55d5276bd5503829be
# source: https://github.com/CTNOriginals/getkeystate/blob/197603a7a57177510262a840f846956a43017ba0/Makefile

MAKE               := make --no-print-directory

DESCRIBE           := $(shell git describe --match "v*" --always --tags)
DESCRIBE_PARTS     := $(subst -, ,$(DESCRIBE))

VERSION_TAG        := $(word 1,$(DESCRIBE_PARTS))
COMMITS_SINCE_TAG  := $(word 2,$(DESCRIBE_PARTS))

VERSION            := $(subst v,,$(VERSION_TAG))
VERSION_PARTS      := $(subst ., ,$(VERSION))

MAJOR              := $(word 1,$(VERSION_PARTS))
MINOR              := $(word 2,$(VERSION_PARTS))
PATCH              := $(word 3,$(VERSION_PARTS))

NEXT_MAJOR         := $(shell echo $$(($(MAJOR)+1)))
NEXT_MINOR         := $(shell echo $$(($(MINOR)+1)))
NEXT_PATCH          = $(shell echo $$(($(PATCH)+$(COMMITS_SINCE_TAG))))

ifeq ($(strip $(COMMITS_SINCE_TAG)),)
CURRENT_VERSION_PATCH := $(MAJOR).$(MINOR).$(PATCH)
CURRENT_VERSION_MINOR := $(CURRENT_VERSION_PATCH)
CURRENT_VERSION_MAJOR := $(CURRENT_VERSION_PATCH)
else
CURRENT_VERSION_PATCH := $(MAJOR).$(MINOR).$(NEXT_PATCH)
CURRENT_VERSION_MINOR := $(MAJOR).$(NEXT_MINOR).0
CURRENT_VERSION_MAJOR := $(NEXT_MAJOR).0.0
endif

# --- Version commands ---
.PHONY: version

version:
	@echo "v$(CURRENT_VERSION_PATCH)"

# -- Project --
.PHONY: run wrun build-win build-linux build

run:
	go run ./main.go

wrun: # requires wgo: https://github.com/bokwoon95/wgow
	wgo run ./main.go

build-win:
	GOOS=windows GOARCH=386 go build -o ./build/BitburnerGoFilesync.exe ./main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/BitburnerGoFilesync ./main.go

build:
	$(MAKE) build-win
	$(MAKE) build-linux

# -- Release --
.PHONY: version-update patch minor major

version-update:
	git tag "v$(VERS)"
	git push --tags

patch:
	$(MAKE) version-update VERS=$(MAJOR).$(MINOR).$(NEXT_PATCH)

minor:
	$(MAKE) version-update VERS=$(MAJOR).$(NEXT_MINOR).0

major:
	$(MAKE) version-update VERS=$(NEXT_MAJOR).0.0