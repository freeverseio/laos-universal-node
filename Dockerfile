# Use an official Go runtime as a parent image
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache ca-certificates git alpine-sdk linux-headers

# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY ./go .

RUN go mod download
# Build the Go program
RUN CGO_ENABLED=1 go build -race -installsuffix 'static' -o app .

# Final stage: the running container.
FROM alpine AS final
RUN apk --no-cache add ca-certificates
# Set the working directory to /app
WORKDIR /app
COPY --from=builder  /app /app

# Run the Go program when the container starts
CMD ["./app"]

# Perform any further action as an unprivileged user.
USER nobody:nobody