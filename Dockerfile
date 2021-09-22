FROM scratch

COPY ./nats-resender /
ENTRYPOINT ["/nats-resender"]
