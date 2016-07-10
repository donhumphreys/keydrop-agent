FROM scratch
COPY docker-keydrop-agent /keydrop-agent
ENTRYPOINT ["/keydrop-agent"]
