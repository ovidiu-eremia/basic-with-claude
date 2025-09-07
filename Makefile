.PHONY: help test loc-prod loc-test loc-history go-files-by-size loc-diff coverage-history

.DEFAULT_GOAL := help

help:
	@echo "Available commands:"
	@echo "  help             Show this help message"
	@echo "  test             Run all tests"
	@echo "  loc-prod         Count production lines of code (excludes tests)"
	@echo "  loc-test         Count test lines of code"
	@echo "  loc-history      Show production LOC by commit in chronological order"
	@echo "  loc-diff         Compare production LOC between current tree and HEAD"
	@echo "  go-files-by-size List all Go files by descending order of size"

test:
	go test ./...

loc-prod:
	@cloc --exclude-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print "Production LOC:", $$5}'

loc-test:
	@cloc --include-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print "Test LOC:", $$5}'

loc-history:
	@./scripts/loc-history.sh

loc-diff:
	@./scripts/loc-diff.sh

go-files-by-size:
	@cloc --by-file --include-ext=go . --quiet

coverage-history:
	@./scripts/coverage-history.sh $(RANGE)
