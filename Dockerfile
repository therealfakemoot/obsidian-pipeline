FROM golang:latest AS builder
ADD . /opt/obp
WORKDIR /opt/obp
# RUN go build -o baas cmd/main.go
RUN make build-alpine

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ARG VERSION=*
COPY --from=builder /opt/obp/dist/obp-$VERSION-alpine_amda64/obp /bin/obp
RUN chmod +x /bin/obp
