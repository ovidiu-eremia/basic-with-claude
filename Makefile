.PHONY: help test loc-prod loc-test loc-history go-files-by-size loc-diff

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

loc-diff:
	@echo "Comparing production Go LOC between current tree and HEAD..."
	@echo ""
	@current_loc=$$(cloc --exclude-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print $$5}' || echo "0"); \
	current_branch=$$(git branch --show-current); \
	git stash push -q -m "temp stash for loc-diff" >/dev/null 2>&1 || true; \
	git checkout HEAD >/dev/null 2>&1; \
	head_loc=$$(cloc --exclude-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print $$5}' || echo "0"); \
	git checkout $$current_branch >/dev/null 2>&1; \
	git stash pop -q >/dev/null 2>&1 || true; \
	if [ -z "$$current_loc" ]; then current_loc=0; fi; \
	if [ -z "$$head_loc" ]; then head_loc=0; fi; \
	diff=$$((current_loc - head_loc)); \
	printf "HEAD LOC:    %6s\n" $$head_loc; \
	printf "Current LOC: %6s\n" $$current_loc; \
	printf "Difference:  %6s" $$diff; \
	if [ $$diff -gt 0 ]; then printf " (+%s lines added)\n" $$diff; \
	elif [ $$diff -lt 0 ]; then printf " (%s lines removed)\n" $$diff; \
	else printf " (no change)\n"; fi

go-files-by-size:
	@cloc --by-file --include-ext=go . --quiet