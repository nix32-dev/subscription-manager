FROM golang:1.25.7-bookworm
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o subscriptions ./cmd
CMD ["./subscriptions"]
