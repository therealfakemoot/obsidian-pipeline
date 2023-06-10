FROM golang:latest AS builder
ADD . /opt/obp
WORKDIR /opt/obp
# RUN go build -o baas cmd/main.go
RUN make alpine-binary

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /opt/obp/bin/obp-alpine /bin/obp
RUN chmod +x /bin/obp
