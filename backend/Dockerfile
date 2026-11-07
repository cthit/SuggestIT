# The binary is built outside Docker: run `make build` first (see Makefile).
FROM scratch

# CA certificates for HTTPS calls to Gamma.
COPY --from=alpine:3.24.4 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY bin/suggestit /suggestit

USER 65534:65534
EXPOSE 3001
ENTRYPOINT ["/suggestit"]
