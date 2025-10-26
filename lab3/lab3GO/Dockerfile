FROM golang:latest

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY tasks/*.go ./tasks/

RUN go build -o main .

CMD ["./main"]



