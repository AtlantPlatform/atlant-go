FROM golang:1.11.5-stretch
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y update \
    && apt-get install -y git gcc libc-dev make

# Get all dependencies for atlant-go and for surgical extract
# Without having them - we are unable to check the build
# we cannot do "go get -u" as some of the deps are broken
RUN  go get -u "github.com/ethereum/go-ethereum/accounts/abi/bind" \
  && go get -u "github.com/ethereum/go-ethereum/common" \
  && go get -u "github.com/ethereum/go-ethereum/ethclient" \
  && go get -u "github.com/ethereum/go-ethereum/rpc" 

RUN go get -u "github.com/elazarl/go-bindata-assetfs" \
 && go get -u "github.com/gin-gonic/gin" \
 && go get -u "github.com/glycerine/go-capnproto"

RUN  go get -u "github.com/jawher/mow.cli" \
  && go get -u "github.com/oklog/ulid" \
  && go get -u "github.com/xlab/catcher" \
  && go get -u "github.com/xlab/closer" \
  && go get -u "github.com/serialx/hashring" \
  && go get -u "github.com/sirupsen/logrus" \
  && go get -u "github.com/stretchr/testify/require"

# - - - - - -
# "Badger" dependency has no backward compatibility,
# So we are taking it from the old branch
# and it has own dependencies:
RUN  go get -u github.com/AndreasBriese/bbloom \
  && go get -u github.com/dgryski/go-farm \ 
  && go get -u github.com/pkg/errors \
  && go get -u golang.org/x/net/trace
WORKDIR /go/src/github.com/dgraph-io
RUN  git config --global advice.detachedHead false \
  && git clone --branch v1.5.4 --single-branch https://github.com/dgraph-io/badger \
  && cd badger \
  && go build .
# - - - - - -

# Now - building Atlant Platform
WORKDIR /go/src/github.com/AtlantPlatform
RUN git clone --depth=1 https://github.com/AtlantPlatform/atlant-go
WORKDIR /go/src/github.com/AtlantPlatform/atlant-go
RUN go get -u .
RUN go build .


# Runner-container (contains only binary to run)
FROM golang:1.11.5-stretch
WORKDIR /AtlantPlatform/
COPY --from=0 /go/src/github.com/AtlantPlatform/atlant-go/atlant-go /AtlantPlatform/
RUN ./atlant-go -T init
# Swarm port to be exposed
EXPOSE 33770
# HTTP API port to be exposed
EXPOSE 33780
# Temporary files folder:
VOLUME /AtlantPlatform/var
CMD ["./atlant-go"]