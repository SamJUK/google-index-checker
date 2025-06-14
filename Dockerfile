FROM alpine:3.22.0

WORKDIR /app

COPY bin/google-index-checker /usr/local/bin/google-index-checker
RUN chmod +x /usr/local/bin/google-index-checker

ENTRYPOINT ["/usr/local/bin/google-index-checker"]
