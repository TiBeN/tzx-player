#!/bin/sh
#
# Project builder

set -e

script_path=$(readlink -f "$(dirname "$0")")

src_path="$script_path"/src
bin_path="$script_path"/bin

cd "$src_path"

# Build
mkdir -p "$bin_path"

# Linux
GOOS=linux GOARCH=amd64 go build -o "$bin_path/tzx-player"

# Windows (requires mingw-w64: sudo pacman -S mingw-w64-gcc)
# Download Windows PortAudio if not already present
win_deps_path="$script_path/win-deps"
if [ ! -d "$win_deps_path/portaudio" ]; then
    echo "Downloading PortAudio source..."
    mkdir -p "$win_deps_path"
    curl -L -o "$win_deps_path/portaudio.tar.gz" \
        "https://github.com/PortAudio/portaudio/archive/refs/tags/v19.7.0.tar.gz"
    tar xzf "$win_deps_path/portaudio.tar.gz" -C "$win_deps_path"
    rm "$win_deps_path/portaudio.tar.gz"

    # Build PortAudio for Windows using mingw
    cd "$win_deps_path/portaudio-19.7.0"
    ./configure --host=x86_64-w64-mingw32 --prefix="$win_deps_path/portaudio" \
        --enable-static --disable-shared
    make
    make install
    cd "$src_path"
fi

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
    PKG_CONFIG_PATH="$win_deps_path/portaudio/lib/pkgconfig" \
    CGO_LDFLAGS="-static" \
    go build -o "$bin_path/tzx-player.exe"