#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(glide nv| tr -s "\n" " "|sed -s "s#/\.\.\.##g"); do
    if ls $d/*.go &> /dev/null; then
        go test -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    fi
done
