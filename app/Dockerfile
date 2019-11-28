FROM golang:alpine AS builder

WORKDIR /go
COPY . .
COPY ./static ./bin/static
RUN apk add --no-cache git
RUN go get -d -v ./...
ARG CGO_ENABLED=0
RUN go build -o ./bin/app ./src/shortenurl/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /go/bin .

CMD ["./app"]
