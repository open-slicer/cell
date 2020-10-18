FROM golang:1.15

WORKDIR /opt/cell
COPY . .

RUN go install -v ./...
