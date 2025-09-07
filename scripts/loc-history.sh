#!/bin/bash
# ABOUTME: Shows production LOC by commit in chronological order
# ABOUTME: Uses a temporary detached git worktree to avoid touching your tree

set -euo pipefail

echo "Commit         Date                Production LOC   Diff  Message"
echo "-------        ----                --------------   ----  -------"

prev_loc=""

# Create a temporary clone rooted at HEAD; stop on failure
wt_dir=$(mktemp -d 2>/dev/null || mktemp -d -t lochistory)
cleanup() {
  rm -rf "$wt_dir" >/dev/null 2>&1 || true
}
trap cleanup EXIT

if ! git clone --quiet . "$wt_dir" >/dev/null 2>&1; then
  echo "Failed to clone repository into temporary directory" >&2
  exit 1
fi

# Iterate commits in chronological order (avoid mapfile for broader compatibility)
for commit in $(git -C "$wt_dir" rev-list --reverse HEAD); do
    # Checkout each commit in the clone
    git -C "$wt_dir" checkout -q "$commit" >/dev/null 2>&1

    date=$(git -C "$wt_dir" show -s --format=%ci "$commit" | cut -d' ' -f1)
    short_commit=$(echo "$commit" | cut -c1-8)

    # Count Go production LOC from within the isolated tree (tolerant to failures)
    set +o pipefail
    loc=$(cd "$wt_dir" && cloc --exclude-content='testing\.T' --include-ext=go . --quiet 2>/dev/null | awk '/^Go/ {print $5}' || true)
    set -o pipefail
    if [ -z "$loc" ]; then loc=0; fi

    if [ -z "$prev_loc" ]; then
        diff=$loc
    else
        diff=$((loc - prev_loc))
    fi

    message=$(git -C "$wt_dir" show -s --format=%s "$commit")
    printf "%-14s %-19s %12d %+8d  %s\n" "$short_commit" "$date" "$loc" "$diff" "$message"
    prev_loc=$loc
done
