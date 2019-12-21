#!/bin/sh

# go get -u github.com/mitchellh/gox

gox                                              \
    -arch="amd64"                                \
    -os="linux darwin windows"                   \
    -output=".build/{{.Dir}}-{{.OS}}-{{.Arch}}" \
    -rebuild                                     \
    .