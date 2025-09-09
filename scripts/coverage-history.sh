#!/bin/bash
# ABOUTME: Shows Go test coverage by commit in chronological order
# ABOUTME: Uses a temporary detached clone to avoid touching your tree
#
# Uses combined coverage (-coverpkg=./...) to include acceptance test coverage.
# Previously used basic coverage which missed functions only tested via acceptance tests.
# Coverage calculation uses 'go tool cover -func' for accuracy vs manual parsing.

set -euo pipefail

echo "Commit         Date        Coverage(%)  Diff(%)  Message"
echo "-------        ----        -----------  -------  -------"

prev_cov=""

# Optional revision range (e.g., "HEAD~50..HEAD"). Defaults to HEAD history.
REV_RANGE="${1:-HEAD}"

# Create a temporary clone rooted at HEAD (or provided range start)
wt_dir=$(mktemp -d 2>/dev/null || mktemp -d -t coverhistory)
cleanup() {
  rm -rf "$wt_dir" >/dev/null 2>&1 || true
}
trap cleanup EXIT
# Allow Ctrl-C to stop immediately (not swallowed by non-fatal blocks)
trap 'exit 130' INT

if ! git clone --quiet . "$wt_dir" >/dev/null 2>&1; then
  echo "Failed to clone repository into temporary directory" >&2
  exit 1
fi

# Compute ordered list of commits for the given range from inside the clone
commits=$(git -C "$wt_dir" rev-list --reverse "$REV_RANGE")

# Iterate commits without using mapfile for portability
for commit in $commits; do
  # Checkout each commit
  git -C "$wt_dir" checkout -q "$commit" >/dev/null 2>&1 || continue

  date=$(git -C "$wt_dir" show -s --format=%ad --date=short "$commit")
  short_commit=$(echo "$commit" | cut -c1-8)
  message=$(git -C "$wt_dir" show -s --format=%s "$commit")

  cov_pct=""
  diff_str=""

  # Run tests with combined coverage; tolerate failures and mark as n/a
  set +e
  (cd "$wt_dir" && go test -count=1 -coverprofile=coverage.out -coverpkg=./... ./... >/dev/null 2>&1)
  test_status=$?
  set -e

  # If interrupted, exit now so Ctrl-C works as expected
  if [ ${test_status:-0} -eq 130 ]; then
    exit 130
  fi

  if [ $test_status -eq 0 ] && [ -f "$wt_dir/coverage.out" ]; then
    # Use go tool cover -func to get accurate coverage percentage
    # Extract total coverage percentage from the last line
    cov_pct=$(cd "$wt_dir" && go tool cover -func=coverage.out 2>/dev/null | tail -1 | sed -n 's/.*[[:space:]]\([0-9][0-9]*\.[0-9][0-9]*\)%$/\1/p')
  fi

  if [ -z "$cov_pct" ]; then
    cov_disp="n/a"
    diff_str="   n/a"
  else
    cov_disp="$cov_pct%"
    if [ -z "$prev_cov" ]; then
      diff_str=$(printf '%+6.2f' "$cov_pct")
    else
      # Compute numeric diff with awk to handle floats precisely
      delta=$(awk -v a="$cov_pct" -v b="$prev_cov" 'BEGIN { printf "%.2f", (a - b) }')
      diff_str=$(printf '%+6.2f' "$delta")
    fi
    prev_cov="$cov_pct"
  fi

  printf "%-14s %-10s %11s %8s  %s\n" "$short_commit" "$date" "$cov_disp" "$diff_str" "$message"
done
