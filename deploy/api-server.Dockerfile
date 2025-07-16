FROM golang:1.24.2

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o api-server ./cmd/api-server

CMD ["./api-server"]
