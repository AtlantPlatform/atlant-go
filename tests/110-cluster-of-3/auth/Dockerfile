FROM atlantplatform/atlant-node:latest
WORKDIR /AtlantPlatform/
COPY auth.conf ./
EXPOSE   33700
ENTRYPOINT ["./atlant-auth", "-s", "./auth.conf", "--log-level", "3" ]
