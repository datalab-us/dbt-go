PROGRAM_NAME := dg
GOPATH_BIN := $(shell go env GOPATH)/bin
UNAME := $(shell uname -s)
PATH_UPDATE_SCRIPT := $(if $(findstring MINGW, $(UNAME)),~/.bash_profile,~/.zshrc)

.PHONY: all build install check-path

all: check-path build install

check-path:
	@if echo "$(PATH)" | grep -q "$(GOPATH_BIN)"; then \
		echo "$(GOPATH_BIN) is already in PATH"; \
	else \
		echo "Adding $(GOPATH_BIN) to PATH"; \
		export PATH="$(GOPATH_BIN):$(PATH)"; \
		echo 'export PATH="$(GOPATH_BIN):$$PATH"' >> $(PATH_UPDATE_SCRIPT); \
		echo "Remember to run 'source $(PATH_UPDATE_SCRIPT)' (or open a new terminal) to update your PATH"; \
	fi

build:
	@echo "Building $(PROGRAM_NAME)..."
	@go build || (echo "Build failed" && exit 1)
	@echo "Build successful"

install:
	@echo "Installing $(PROGRAM_NAME)..."
	@if [ "$(UNAME)" = "Linux" ] || [ "$(UNAME)" = "Darwin" ]; then \
		sudo mv $(PROGRAM_NAME) /usr/local/bin/ || (echo "Installation failed" && exit 1); \
	elif echo "$(UNAME)" | grep -q MINGW; then \
		mv $(PROGRAM_NAME).exe $(GOPATH_BIN) || (echo "Installation failed" && exit 1); \
	fi
	@echo "Installation successful."
	@echo "Ensure you run '$(PROGRAM_NAME)' from the root of your dbt project."

release:
	@echo "Releasing $(PROGRAM_NAME)..."
	@version=$(shell grep 'Version =' ./version/version.go | awk '{print $$NF}' | tr -d '"') && \
	git tag -a $$version -m "Release $$version" && \
	git push origin main --follow-tags
	@echo "Release successful."