#!/bin/bash
# ABOUTME: Shows Go test coverage percentage by commit in chronological order
# ABOUTME: Uses a temporary git worktree to avoid touching the current tree

set -euo pipefail

echo "Commit         Date                Coverage %   Diff   Message"
echo "-------        ----                ----------   ----   -------"

repo_root=$(git rev-parse --show-toplevel)
current_branch=$(git branch --show-current || true)
prev_cov=""

# Prepare a temporary worktree at HEAD
wt_dir=$(mktemp -d 2>/dev/null || mktemp -d -t covwt)
cleanup() {
  git worktree remove --force "$wt_dir" >/dev/null 2>&1 || true
}
trap cleanup EXIT

git worktree add --detach -q "$wt_dir" HEAD >/dev/null 2>&1

# Collect commits in chronological order
mapfile -t commits < <(git rev-list --reverse HEAD)

for commit in "${commits[@]}"; do
    # Checkout commit inside the worktree
    git -C "$wt_dir" checkout -q "$commit" >/dev/null 2>&1

    date=$(git show -s --format=%ci "$commit" | cut -d' ' -f1)
    short_commit=$(echo "$commit" | cut -c1-8)

    # Create a temp coverage profile and run tests from the worktree
    tmpfile=$(mktemp 2>/dev/null || mktemp -t coverprofile)
    if (cd "$wt_dir" && go test -count=1 -coverpkg=./... ./... -coverprofile="$tmpfile" -covermode=atomic >/dev/null 2>&1); then
        cov=$(go tool cover -func="$tmpfile" 2>/dev/null | awk '/^total:/ {print $3}' | tr -d '%')
    else
        cov=""
    fi
    rm -f "$tmpfile" 2>/dev/null || true

    # Fallback if coverage not available
    if [ -z "$cov" ]; then cov=0; fi

    # Compute diff vs previous (two decimals)
    if [ -z "$prev_cov" ]; then
        diff=$cov
    else
        diff=$(awk -v a="$cov" -v b="$prev_cov" 'BEGIN{printf "%.2f", (a - b)}')
    fi

    message=$(git show -s --format=%s "$commit")

    # Print: commit, date, coverage (2 decimals), diff (2 decimals with sign), message
    printf "%-14s %-19s %10.2f %% %+6.2f  %s\n" "$short_commit" "$date" "$cov" "$diff" "$message"

    prev_cov=$cov
done
