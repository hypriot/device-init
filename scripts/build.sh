#!/usr/bin/env bash
set -e

OSARCH="linux/arm linux/amd64"
WORKDIR="/opt/gopath/src/github.com/hypriot/device-init"

cd $WORKDIR
GOPATH="`godep path`:$GOPATH" gox -osarch="${OSARCH}"

# ensure that the travis-ci user can access the binaries
chmod a+rw device-init_linux*
