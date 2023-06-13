FROM golang:latest AS builder
ADD . /opt/obp
WORKDIR /opt/obp
RUN make build-alpine

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ARG VERSION=*
COPY --from=builder /opt/obp/dist/obp-*-alpine_amd64/obp /bin/obp
RUN chmod +x /bin/obp
