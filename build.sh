#!/bin/sh
#
# Project builder

set -e

script_path=$(readlink -f "$(dirname "$0")")

src_path="$script_path"/src
bin_path="$script_path"/bin

export GOPATH="$src_path/.go"
export GOCACHE="$src_path/.go/cache"

cd "$src_path"

# Fetch external packages

go get -t ./...

# Build
mkdir -p "$bin_path"
go build -o "$bin_path/tzx-player"