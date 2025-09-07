#!/bin/bash
# ABOUTME: Shows production LOC by commit in chronological order
# ABOUTME: Iterates through git history and counts production Go lines for each commit

echo "Commit         Date                Production LOC   Diff  Message"
echo "-------        ----                --------------   ----  -------"

current_branch=$(git branch --show-current)
prev_loc=""

git rev-list --reverse HEAD | while read commit; do
    git checkout $commit >/dev/null 2>&1
    date=$(git show -s --format=%ci $commit | cut -d' ' -f1)
    short_commit=$(echo $commit | cut -c1-8)
    loc=$(cloc --exclude-content='testing\.T' --include-ext=go . --quiet 2>/dev/null | awk '/^Go/ {print $5}' || echo "0")
    if [ -z "$loc" ]; then loc=0; fi

    if [ -z "$prev_loc" ]; then
        diff=$loc
    else
        diff=$((loc - prev_loc))
    fi

    message=$(git show -s --format=%s $commit)
    printf "%-14s %-19s %12d %+8d  %s\n" $short_commit $date $loc $diff "$message"
    prev_loc=$loc
done

git checkout $current_branch >/dev/null 2>&1
