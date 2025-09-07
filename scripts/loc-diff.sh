#!/bin/bash
# ABOUTME: Compares production Go LOC between current tree and HEAD
# ABOUTME: Shows detailed diff with stashing to handle uncommitted changes safely

echo "Comparing production Go LOC between current tree and HEAD..."
echo ""

current_loc=$(cloc --exclude-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print $5}' || echo "0")
current_branch=$(git branch --show-current)

git stash push -q -m "temp stash for loc-diff" >/dev/null 2>&1 || true
git checkout HEAD >/dev/null 2>&1
head_loc=$(cloc --exclude-content='testing\.T' --include-ext=go . --quiet | awk '/^Go/ {print $5}' || echo "0")
git checkout $current_branch >/dev/null 2>&1
git stash pop -q >/dev/null 2>&1 || true

if [ -z "$current_loc" ]; then current_loc=0; fi
if [ -z "$head_loc" ]; then head_loc=0; fi

diff=$((current_loc - head_loc))

printf "HEAD LOC:    %6s\n" $head_loc
printf "Current LOC: %6s\n" $current_loc
printf "Difference:  %6s" $diff

if [ $diff -gt 0 ]; then 
    printf " (+%s lines added)\n" $diff
elif [ $diff -lt 0 ]; then 
    printf " (%s lines removed)\n" $diff
else 
    printf " (no change)\n"
fi