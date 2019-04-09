#!/usr/bin/env bash
set -e
cd "$(dirname "$0")" && pwd

if ! [ "`which docker`" ]; then
  echo 'Error: docker must be installed' >&2
  exit 1
fi

docker build -t atlantplatform ../../