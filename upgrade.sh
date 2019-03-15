#!/bin/bash

# exit if any command fails
set -e

echo "== Re-creating vendor folder"
# remove vendor from gitignore. it will be included in commits
sed -i '/vendor/d' ./.gitignore
# remove all vendors
rm -rf ./vendor/ go.mod go.sum

echo "== Downloading dependencies"
export GO111MODULE=on && go mod init && go mod vendor && go mod download && export GO111MODULE=off

# getting IPFS done right is still tricky.
# so the tactics is to download it separately,
# and then merge it with our vendor to prevent type mismatches
echo "== Overwriting master branch from go-ipfs"
rm -rf ./vendor/github.com/ipfs/go-ipfs
git clone -q --depth=1 -b master --single-branch https://github.com/ipfs/go-ipfs ./vendor/github.com/ipfs/go-ipfs
rm -rf ./vendor/github.com/ipfs/go-ipfs/.git
cd ./vendor/github.com/ipfs/go-ipfs
echo "== Downloading go-ipfs deps"
export GO111MODULE=on && go mod vendor && go mod download && export GO111MODULE=off
cd -

echo "== Merging go-ipfs deps"
cp -rf ./vendor/github.com/ipfs/go-ipfs/vendor ./
rm -rf ./vendor/github.com/ipfs/go-ipfs/vendor

echo "== Outstanding vendors"
rm -rf ./vendor/github.com/davidlazar/go-crypto
git clone -q --depth=1 https://github.com/davidlazar/go-crypto \
    ./vendor/github.com/davidlazar/go-crypto

echo "== Removing Extra(s)"
cd vendor && find . -type d -name ".git" -exec rm -rf {} + && cd -
rm -rf ./vendor/github.com/ethereum ./vendor/modules.txt go.mod go.sum

# Apply all patches from patches folder
for f in ./patches/*.patch;
do
  echo "== Applying patch $f"
  git apply $f
done;


echo "== Building Atlant-go"
go build

echo "== Atlant-go: Starting, init node"
rm -rf var/ && git checkout var/genesis.json
./atlant-go -T init
./atlant-go
# if it is running and responding, then feel free to commit

