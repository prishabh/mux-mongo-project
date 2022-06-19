FROM golang:1.16-buster AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./

RUN go build -o /mux-mongo-project
EXPOSE 8080
ENTRYPOINT ["/mux-mongo-project"]
