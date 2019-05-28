#!/usr/bin/env bash
set -e

CURRDIR=${PWD##*/} 
if [ "$CURRDIR" != "atlant-go" ]; then
	echo "This script should be launched from the project folder"
  exit 1
fi

# Get the version from the environment, or try to figure it out.
if [ -z $VERSION ]; then
	VERSION=$(awk -F\" '/Version =/ { print $2; exit }' < version/version.go)
fi
if [ -z "$VERSION" ]; then
    echo "Please specify a version."
    exit 1
fi
echo "==> Building version $VERSION..."

# Delete the old dir
echo "==> Removing old directory..."
rm -rf build/pkg
mkdir -p build/pkg

# Get the git commit
GIT_COMMIT="$(git rev-parse --short=8 HEAD)"
GIT_IMPORT="github.com/AtlantPlatform/atlant-go/version"

# Determine the arch/os combos we're building for
XC_TARGETS="darwin/amd64,linux/amd64,linux/arm,linux/386,windows/386,windows/amd64" 
# Build!
# ldflags: -s Omit the symbol table and debug information.
#	         -w Omit the DWARF symbol table.
echo "==> Building..."
GO111MODULE=off
xgo -ldflags="-s -w -X ${GIT_IMPORT}.GitCommit=${GIT_COMMIT}" \
	-targets=${XC_TARGETS} -out "" -dest="build/pkg/" \
	./

# Zip all the files.
echo "==> Packaging..."
for BINARY in $(find ./build/pkg -mindepth 1 -maxdepth 1 -type f); do
	echo "--> ${BINARY}"
	BINARY=$(basename "$BINARY")
	zip "./build/pkg/${BINARY:1}.zip" "./build/pkg/${BINARY}"
	rm "./build/pkg/${BINARY}"
done

# Add "atlant-go" and $VERSION prefix to package name.
rm -rf ./build/dist
mkdir -p ./build/dist
for FILENAME in $(find ./build/pkg -mindepth 1 -maxdepth 1 -type f); do
  FILENAME=$(basename "$FILENAME")
	cp "./build/pkg/${FILENAME}" "./build/dist/atlant-go_${VERSION}_${FILENAME}"
done

# Make the checksums.
pushd ./build/dist
shasum -a256 ./* > "./atlant-go_${VERSION}_SHA256SUMS"
popd

# Done
echo
echo "==> Results:"
ls -hl ./build/dist

exit 0
