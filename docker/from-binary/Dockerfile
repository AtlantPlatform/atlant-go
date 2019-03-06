FROM ubuntu:latest
## NOTE: the binary from releases doesn't work for alpine!!!

RUN apt-get -y update && apt-get install -y unzip
WORKDIR /AtlantPlatform
ADD https://github.com/AtlantPlatform/atlant-go/releases/download/v1.0.0-rc2/atlant-go-linux-amd64.zip ./
RUN unzip atlant-go-linux-amd64.zip && rm -f atlant-go-linux-amd64.zip
RUN mv -f build/pkg/atlant-go-linux-amd64 /AtlantPlatform/atlant-go
RUN ./atlant-go -T init

# Swarm port to be exposed
EXPOSE 33770
# HTTP API port to be exposed
EXPOSE 33780
# Temporary files folder:
VOLUME /AtlantPlatform/var

CMD ["./atlant-go"]