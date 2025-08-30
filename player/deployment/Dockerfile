#-------------------------------------------------------
FROM golang:alpine AS backend-builder
# Add gcc for cgo
RUN apk add build-base
# Add lua5.3
RUN apk add lua5.3 lua5.3-dev
# Fix build issue with lua
RUN (cd /usr/lib && ln -s lua5.3/liblua.so liblua5.3.so)

COPY ./rgcore /rgcore

COPY ./player/go.mod  ./player/go.sum /player/

WORKDIR /player

RUN go mod download
RUN go mod verify

COPY ./player/internal internal
COPY ./player/cmd cmd

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o cmd/main cmd/main.go

#-------------------------------------------------------
ARG PORT

FROM alpine

# Add lua5.3
RUN apk add lua5.3 lua5.3-dev

COPY --from=backend-builder ./player/cmd/main ./

EXPOSE $PORT

ENTRYPOINT ["./main"]

CMD ["-p", $PORT]