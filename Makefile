.PHONY: help test loc-prod loc-test loc-history

.DEFAULT_GOAL := help

help:
	@echo "Available commands:"
	@echo "  help         Show this help message"
	@echo "  test         Run all tests"
	@echo "  loc-prod     Count production lines of code (excludes tests)"
	@echo "  loc-test     Count test lines of code"
	@echo "  loc-history  Show production LOC by commit in chronological order"

test:
	go test ./...

loc-prod:
	@cloc --exclude-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print "Production LOC:", $$5}'

loc-test:
	@cloc --include-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print "Test LOC:", $$5}'

loc-history:
	@echo "Commit         Date                Production LOC"
	@echo "-------        ----                --------------"
	@current_branch=$$(git branch --show-current); \
	git rev-list --reverse HEAD | while read commit; do \
		git checkout $$commit >/dev/null 2>&1; \
		date=$$(git show -s --format=%ci $$commit | cut -d' ' -f1); \
		short_commit=$$(echo $$commit | cut -c1-8); \
		loc=$$(cloc --exclude-content='testing\.T' --include-ext=go . --quiet 2>/dev/null | awk '/^Go/ {print $$5}' || echo "0"); \
		printf "%-14s %-19s %s\n" $$short_commit $$date $$loc; \
	done; \
	git checkout $$current_branch >/dev/null 2>&1