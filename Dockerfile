FROM golang:1.21-alpine3.18 AS builder
# Add C compiler to build the executable with race detection enabled
RUN apk add --no-cache build-base=~0.5

RUN addgroup -g 1000 nodegroup && \
    adduser -D nodeuser -u 1000 -G nodegroup

WORKDIR /app
COPY ./go .
RUN go build -race -o app .

FROM alpine:3.18.4 AS final

WORKDIR /app
COPY --from=builder /app /app

CMD ["./app"]
USER nobody:nobody