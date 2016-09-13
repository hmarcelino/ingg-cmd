#!/bin/bash

set +x

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH="$CURRENT_DIR"

# build for the current os
echo "Building ingg command line"
go build ingg

if [[ "$?" != 0 ]]; then
    echo "Error building ingg command line"
    exit 1
fi

# Build for windows
echo "Building ingg command line for windows"
GOOS=windows GOARCH=amd64 go build ingg

mkdir -p $CURRENT_DIR/dist
mv $CURRENT_DIR/ingg $CURRENT_DIR/ingg.exe $CURRENT_DIR/dist/

cp $CURRENT_DIR/dist/ingg ~/dotfiles/bin/



