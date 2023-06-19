# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:latest AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# This definitely works
# RUN CGO_ENABLED=0 GOOS=linux go build -o obp ./cmd/obp/
RUN make build-alpine


# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /

ARG VERSION=version
COPY --from=build-stage /app/dist/*-alpine/obp /bin/
RUN chmod +x /bin/obp

ARG USER=default
ENV HOME /home/$USER

RUN apk add --update sudo

RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME
