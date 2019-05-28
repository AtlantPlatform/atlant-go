FROM atlantplatform/atlant-node:latest
WORKDIR /AtlantPlatform/
ENV AN_LOG_LEVEL=5

RUN ./atlant-go -T init
EXPOSE 33770 33780

CMD ["./atlant-go"]