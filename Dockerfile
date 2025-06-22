FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chores ./cmd/chores/main.go

# 7. Run the binary
CMD ["./chores"]
