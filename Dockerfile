FROM golang:1.16-alpine as build
WORKDIR /src
ADD . /src
RUN apk --no-cache add build-base git
RUN make

FROM scratch

COPY --from=build /src/build/nats-resender /
ENTRYPOINT ["/nats-resender"]
