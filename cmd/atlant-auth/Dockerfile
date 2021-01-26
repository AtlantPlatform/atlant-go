# Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
# Use of this source code is governed by BSD 3-Clause "New" or "Revised"
# License (BSD 3) that can be found in the LICENSE file.

# This docker file contains binary files together with source code 

FROM golang:1.12-stretch
WORKDIR /go/src/github.com/AtlantPlatform/atlant-auth
COPY . .
ENV GO111MODULE=on 
RUN go mod vendor
RUN go build && go test -v

# Container that contains only binaries
FROM debian:stretch
WORKDIR /AtlantPlatform/
COPY --from=0 /go/src/github.com/AtlantPlatform/atlant-auth/atlant-* /AtlantPlatform/
EXPOSE 33700
CMD ["./atlant-auth"]