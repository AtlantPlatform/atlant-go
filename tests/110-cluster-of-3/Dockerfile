FROM atlantplatform:latest
WORKDIR /AtlantPlatform/

RUN ./atlant-go -T init
EXPOSE 33770 33780

VOLUME /AtlantPlatform/var/fs
VOLUME /AtlantPlatform/var/log
VOLUME /AtlantPlatform/var/state

CMD ["./atlant-go"]