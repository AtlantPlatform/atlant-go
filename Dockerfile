# Copyright 2017-2019 Tensigma Ltd. All rights reserved.
# Use of this source code is governed by Microsoft Reference Source
# License (MS-RSL) that can be found in the LICENSE file.

# This docker file contains binary files together with source code 
# of atlant-go

FROM golang:1.12-stretch
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y update \
    && apt-get install -y git gcc libc-dev make

# Get ETHEREUM libraries - they are not included
RUN  go get -u "github.com/ethereum/go-ethereum/accounts/abi/bind" \
  && go get -u "github.com/ethereum/go-ethereum/common" \
  && go get -u "github.com/ethereum/go-ethereum/ethclient" \
  && go get -u "github.com/ethereum/go-ethereum/rpc"
 
# Now - building Atlant Platform
WORKDIR /go/src/github.com/AtlantPlatform/atlant-go
ENV GO111MODULE=off
COPY . .
RUN go build

# Build Atlant Lite
RUN cd cmd/atlant-lite && go build -o ./../../atlant-lite main.go && cd -
# Build Atlant Keygen
RUN cd cmd/atlant-keygen && go build -o ./../../atlant-keygen main.go && cd -

# Container that contains only binaries
FROM debian:stretch
WORKDIR /AtlantPlatform/
COPY --from=0 /go/src/github.com/AtlantPlatform/atlant-go/atlant-* /AtlantPlatform/
CMD ["./atlant-go"]