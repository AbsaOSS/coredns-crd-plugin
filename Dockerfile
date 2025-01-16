FROM alpine:3.21.2
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
