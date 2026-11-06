# Makefile for labsched.
#
# This file is the CI. GitHub Actions will only ever run `make ci`, so a
# green `make ci` on your laptop means a green build on GitHub.
#
# Make basics:
# - Each "target:" line names a command you can run with `make <target>`.
# - Indented lines under it are the recipe. They MUST start with a tab.
# - Each recipe line runs in its own fresh shell.
# - Make prints each command before running it. A leading @ hides that.
# - Make expands $ itself first, so shell variables need $$.

# The binary name and output folder. `:=` assigns once, immediately.
BINARY := labsched
BIN_DIR := bin

# Running plain `make` runs this target.
.DEFAULT_GOAL := ci

# Make normally treats a target as a file to build, and skips the recipe if
# a file with that name exists. .PHONY says these are commands, not files.
.PHONY: test lint fmt-check tidy-check build ci clean

# Run every test in the module.
#   -race       turns on the race detector. It catches data races, but only
#               on code paths the tests actually run.
#   -shuffle=on runs tests in random order, which exposes tests that depend
#               on each other. It prints the seed so a failure can be rerun.
#   ./...       means this folder and every package below it.
test:
	go test -race -shuffle=on ./...

# Run the linters. With no .golangci.yml yet, golangci-lint uses its
# defaults: govet, staticcheck, errcheck, ineffassign and unused.
lint:
	golangci-lint run

# Fail if any file isn't gofmt-formatted. This checks without changing files.
#   gofmt -l .   lists unformatted files, and prints nothing if all are fine.
#   The backslashes join these lines into one shell command, because each
#   recipe line would otherwise run in its own shell and lose the variable.
#   $$ becomes a single $ by the time the shell sees it.
fmt-check:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

# Fail if go.mod or go.sum don't match what the code imports.
#   -diff prints what `go mod tidy` would change, and exits non-zero if
#   anything would change. It never edits the files.
tidy-check:
	go mod tidy -diff

# Compile the binary into bin/.
#   -o sets the output path. The bin/ folder is created if it's missing.
#   .  builds the main package in this folder.
build:
	go build -o $(BIN_DIR)/$(BINARY) .

# Everything CI runs. The targets after the colon are prerequisites, and Make
# runs them left to right. There's no recipe, because the prerequisites do
# the work. If any one fails, Make stops and the whole target fails.
# The order is cheap and fast checks first, so failures show up quickly.
ci: fmt-check tidy-check lint test build

# Remove build output.
clean:
	rm -rf $(BIN_DIR)
