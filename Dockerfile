FROM golang:1-alpine AS build
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on
RUN apk add --no-cache git
RUN mkdir /out && \
    go test -v ./... && \
    go build -v -o /out/labrador main.go
FROM alpine AS base
COPY --from=build /out/labrador /labrador
RUN adduser --disabled-password --no-create-home labrador
USER labrador
CMD ["/labrador"]