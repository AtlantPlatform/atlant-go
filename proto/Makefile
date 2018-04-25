# go get -u github.com/glycerine/go-capnproto
# go get -u github.com/glycerine/go-capnproto/capnpc-go

all: fmt
	capnp compile -I $(GOPATH)/src/github.com/glycerine/go-capnproto -ogo fs.capnp rs.capnp

fmt:
	./capfmt fs.capnp
	./capfmt rs.capnp
