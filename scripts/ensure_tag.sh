#!/usr/bin/env bash

version=${version:-"v$(gsemver bump)"}

if [[ "$version" != "v" && -z "`git tag -l $version`" ]];then
    git tag -a -m "release version $version" $version
fi
