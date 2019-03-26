FROM atlantplatform:latest

WORKDIR /AtlantPlatform/
RUN echo "DO NOT CACHE" > do-not-cache.node3

RUN ./atlant-go -T --testnet-key d69495233081d310a5bc4b53e8212021c2b8564ec3d70887727fe70be24b15a8 init 
EXPOSE 33770 33780

# VOLUME /AtlantPlatform/var/fs
# VOLUME /AtlantPlatform/var/log
# VOLUME /AtlantPlatform/var/state

CMD ["./atlant-go"]