FROM golang:1.23

# Set destination for COPY
WORKDIR /app

# Download Go modules
# COPY go.mod .
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN go build -o ./chores ./cmd/chores


# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/chores"]
