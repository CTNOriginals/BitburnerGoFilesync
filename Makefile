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

.DEFAULT_GOAL := help
# --- Version commands ---
.PHONY: help version proto

help: ##@help Display all commands and descriptions
	@awk 'BEGIN {FS = ":.*##@"; printf "\nUsage:\n  make <target>\n"} \
	/^[.a-zA-Z_-]+:.*?##@/ { \
		split($$2, parts, " "); \
		section = parts[1]; \
		description = substr($$2, length(section) + 2); \
		sections[section] = sections[section] sprintf("  \033[36m%-15s\033[0m %s\n", $$1, description); \
	} \
	END { \
		for (section in sections) { \
			printf "\n\033[1m%s\033[0m\n", section; \
			printf "%s", sections[section]; \
		} \
	}' $(MAKEFILE_LIST)

list: ##@help List all targets and their commands
	@awk 'BEGIN { \
		target = ""; cmds = ""; \
	} \
	/^[.a-zA-Z_-]+:/ && !/^\.PHONY/ { \
		if (target != "" && cmds != "") { \
			printf "  \033[36m%-15s\033[0m\n%s\n", target, cmds; \
		} \
		split($$0, a, ":"); \
		target = (a[1] == "help" || a[1] == "list") ? "" : a[1]; \
		cmds = ""; \
	} \
	/^\t/ && target != "" { \
		cmds = cmds "    " substr($$0, 2) "\n"; \
	} \
	END { \
		if (target != "" && cmds != "") { \
			printf "  \033[36m%-15s\033[0m\n%s\n", target, cmds; \
		} \
	}' $(MAKEFILE_LIST)

version: ##@help Log the current version
	@echo "v$(CURRENT_VERSION_PATCH)"

# proto: ##@help For prototyping makefile functionality
# 	@echo hello $1

# -- Git --
.PHONY: git-graph adog

git-graph: ##@git Log decorated graph
	git log --all --decorate --oneline --graph
	# git log --graph --abbrev-commit --decorate --format=format:'%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(auto)%d%C(reset)' --all

# -- Run --
.PHONY: run debug test test-coverage wrun

WGO_INCLUDE := -file .go -file .toml -file Makefile

run: ##@run Run normally. Pass arguments like so: args="arg1 arg2 ...".
	go run . $(args)

debug: ##@run Run with --test $(testargs).
	go run . $(args) --test $(testargs)

test: ##@run go test $(args); for all packaged that contain at least 1 *_test.go script.
	go test $(args) $$(go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

test-coverage: ##@run Generate an preview test coverage in html form.
	$(MAKE) test args="-coverprofile tmp/cover.out"
	go tool cover -html="tmp/cover.out" -o "tmp/cover.html"
	xdg-open tmp/cover.html

wrun: ##@run Run a make target and restart on file change. make wrun <wgoargs="args..."> target=[TARGET]. Requires wgo: https://github.com/bokwoon95/wgo
	wgo $(WGO_INCLUDE) $(wgoargs) $(MAKE) $(target)
wrunstdin: ##@run Run a make wrun with wgoargs -stdin.
	$(MAKE) wrun wgoargs="-stdin $(wgoargs)" target="$(target)"

# -- Build --
.PHONY: build-win build-linux build-mac build

build-win: ##@build Build for windows. Binary will be located at ./build/
	GOOS=windows GOARCH=amd64 go build -o ./build/BitburnerGoFilesync_win.exe .

build-linux: ##@build Build for linux. Binary will be located at ./build/
	GOOS=linux GOARCH=amd64 go build -o ./build/BitburnerGoFilesync_linux .

build-mac: ##@build Build for linux. Binary will be located at ./build/
	GOOS=darwin GOARCH=amd64 go build -o ./build/BitburnerGoFilesync_mac .

build: ##@build Build for both windows and linux. Binary will be located at ./build/
	$(MAKE) build-win
	$(MAKE) build-linux
	$(MAKE) build-mac

# -- Release --
.PHONY: tag patch minor major

tag: ##@versioning Push tags
	git push --tags

patch: ##@versioning Create and add a patch tag (vx.x.+commits)
	git tag "v$(MAJOR).$(MINOR).$(NEXT_PATCH)"

minor: ##@versioning Create and add a minor tag (vx.+1.x)
	git tag "v$(MAJOR).$(NEXT_MINOR).0"

major: ##@versioning Create and add a major tag (v+1.x.x)
	git tag "v$(NEXT_MAJOR).0.0"
