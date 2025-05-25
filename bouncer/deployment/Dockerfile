FROM golang:alpine AS backend-builder
# Add gcc for cgo
RUN apk add build-base

COPY ./rgcore/rgentities /rgcore/rgentities

COPY ./bouncer /bouncer

WORKDIR /bouncer

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./main ./cmd/main.go

EXPOSE 5555

ENTRYPOINT ["./main"]