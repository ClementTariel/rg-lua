FROM golang:alpine AS backend-builder
# Add gcc for cgo
RUN apk add build-base
# Add lua5.3
RUN apk add lua5.3 lua5.3-dev
# Fix build issue with lua
RUN (cd /usr/bin && ln -s lua5.3 lua)
RUN (cd /usr/lib && ln -s lua5.3/liblua.so liblua5.3.so)

COPY ./rgcore /rgcore
COPY ./matchmaker /matchmaker

WORKDIR /matchmaker

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./main ./cmd/main.go

EXPOSE 4444

ENTRYPOINT ["./main"]