FROM golang:alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /new-world-dev-tracker

ENTRYPOINT echo "Starting New World Dev Tracker"
ENTRYPOINT [ "/new-world-dev-tracker" ]
