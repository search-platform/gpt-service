FROM golang:1.21 AS builder

ARG TARGETARCH
ARG PROJECT

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags "-s -w" -o /go/bin/app ./cmd/$PROJECT

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/app .

EXPOSE 8080

CMD ["./app"]
