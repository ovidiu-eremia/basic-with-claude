#!/bin/bash
# ABOUTME: Shows Go test coverage by commit in chronological order
# ABOUTME: Uses a temporary detached clone to avoid touching your tree

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

  # Run tests with coverage; tolerate failures and mark as n/a
  set +e
  (cd "$wt_dir" && go test -count=1 ./... -coverprofile=coverage.out >/dev/null 2>&1)
  test_status=$?
  set -e

  # If interrupted, exit now so Ctrl-C works as expected
  if [ ${test_status:-0} -eq 130 ]; then
    exit 130
  fi

  if [ $test_status -eq 0 ] && [ -f "$wt_dir/coverage.out" ]; then
    # Compute coverage directly from coverage.out to avoid go tool cover pitfalls
    # Format: mode: set; then lines like: file:line1,col1,line2,col2 statements count
    cov_pct=$(awk '
      BEGIN { total=0; covered=0 }
      NR==1 { next } # skip mode line
      {
        stmts=$(NF-1); cnt=$NF;
        total+=stmts; if (cnt>0) covered+=stmts;
      }
      END {
        if (total>0) printf "%.2f", (covered/total*100);
      }
    ' "$wt_dir/coverage.out")
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
