FROM golang:alpine AS backend-builder
# Add gcc for cgo
RUN apk add build-base

COPY ./rgcore/rgconst /rgcore/rgconst
COPY ./rgcore/rgentities /rgcore/rgentities
COPY ./rgcore/rgutils /rgcore/rgutils

COPY ./referee /referee

WORKDIR /referee

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./main ./cmd/main.go

EXPOSE 3333

ENTRYPOINT ["./main"]