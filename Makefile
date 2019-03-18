# Copyright 2017-2019 Tensigma Ltd. All rights reserved.
# Use of this source code is governed by Microsoft Reference Source
# License (MS-RSL) that can be found in the LICENSE file.

all:

clean:
	rm -rf var/fs var/state

clean-cluster2: clean-node1 clean-node2
init-cluster2: init-node1 init-node2

clean-cluster3: clean-node1 clean-node2 clean-node3
init-cluster3: init-node1 init-node2 init-node3

init-node1:
	atlant-go -T -F var/fs1 -S var/state1 init

init-node2:
	atlant-go -T -F var/fs2 -S var/state2 init

init-node3:
	atlant-go -T -F var/fs3 -S var/state3 init

clean-node1:
	rm -rf var/fs1 var/state1

clean-node2:
	rm -rf var/fs2 var/state2

clean-node3:
	rm -rf var/fs3 var/state3

run-node1: export IPFS_LOGGING = warning
run-node1:
	atlant-go -E 0x0  -F var/fs1 -S var/state1 -L ":33771" -W ":33781" -l 5 $(COMMAND)

run-node2: export IPFS_LOGGING = warning
run-node2:
	atlant-go -E 0x0 -F var/fs2 -S var/state2 -L ":33772" -W ":33782" -l 5 $(COMMAND)

run-node3: export IPFS_LOGGING = warning
run-node3:
	atlant-go -E 0x0 -F var/fs3 -S var/state3 -L ":33773" -W ":33783" -l 5 $(COMMAND)

test-node1-api:
	go test -short -E 0x0 -F var/fs1 -S var/state1 -L :33771 -W 0.0.0.0:33781

test-node2-api:
	go test -short -E 0x0 -F var/fs2 -S var/state2 -L :33772 -W 0.0.0.0:33782

test-node3-api:
	go test -short -E 0x0 -F var/fs3 -S var/state3 -L :33773 -W 0.0.0.0:33783

install:
	go install -tags testing github.com/AtlantPlatform/atlant-go

GOTOOLS = \
	github.com/golang/dep/cmd/dep \
	gopkg.in/alecthomas/gometalinter.v2

PACKAGES=$(shell go list ./... | grep -v '/vendor/')

BUILD_FLAGS = -ldflags "-X github.com/AtlantPlatform/atlant-go/version.GitCommit=`git rev-parse --short=8 HEAD`"

### Distribution

# dist builds binaries for all platforms and packages them for distribution
dist:
	sh -c "'$(CURDIR)/deploy/dist.sh'"

### Tools & dependencies

check-tools:
	@# https://stackoverflow.com/a/25668869
	@echo "Found tools: $(foreach tool,$(notdir $(GOTOOLS)),\
        $(if $(shell which $(tool)),$(tool),$(error "No $(tool) in PATH")))"

get-tools:
	@echo "--> Installing tools"
	go get -u -v $(GOTOOLS)
	@gometalinter.v2 --install

update-tools:
	@echo "--> Updating tools"
	@go get -u $(GOTOOLS)

#Run this from CI
get-vendor-deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure -vendor-only

#Run this locally.
ensure-deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/AtlantPlatform/atlant-go -d 3 | dot -Tpng -o dependency-graph.png

### Formatting, linting, and vetting

metalinter:
	@echo "--> Running linter"
	@gometalinter.v2 --vendor --deadline=600s --disable-all  \
		--enable=deadcode \
		--enable=gosimple \
	 	--enable=misspell \
		--enable=safesql \
		./...
		#--enable=gas \
		#--enable=maligned \
		#--enable=dupl \
		#--enable=errcheck \
		#--enable=goconst \
		#--enable=gocyclo \
		#--enable=goimports \
		#--enable=golint \ <== comments on anything exported
		#--enable=gotype \
	 	#--enable=ineffassign \
	   	#--enable=interfacer \
	   	#--enable=megacheck \
	   	#--enable=staticcheck \
	   	#--enable=structcheck \
	   	#--enable=unconvert \
	   	#--enable=unparam \
		#--enable=unused \
	   	#--enable=varcheck \
		#--enable=vet \
		#--enable=vetshadow \

metalinter-all:
	@echo "--> Running linter (all)"
	gometalinter.v2 --vendor --deadline=600s --enable-all --disable=lll ./...

