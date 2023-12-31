FROM golang:1.21-alpine3.18 AS builder
# Add C compiler to build the executable with race detection enabled
RUN apk add --no-cache build-base=~0.5

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
COPY internal internal
COPY cmd cmd

ARG VERSION

WORKDIR /app/cmd
RUN go build -race -ldflags "-X main.version=$VERSION" -o universalnode .

FROM alpine:3.18.4 AS final

WORKDIR /app
COPY --from=builder /app/cmd/universalnode .

RUN chown nobody:nobody .
USER nobody:nobody
ENTRYPOINT ["./universalnode", "-storage_path", ".universalnode"]
