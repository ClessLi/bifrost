#!/usr/bin/env bash

if [ "$1" == '-t' ]; then
    echo "check failure" >&2
    return 2
fi
echo pass