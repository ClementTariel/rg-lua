#-------------------------------------------------------
FROM golang:alpine AS backend-builder
# Add gcc for cgo
RUN apk add build-base

COPY ./rgcore/rgconst /rgcore/rgconst
COPY ./rgcore/rgentities /rgcore/rgentities
COPY ./rgcore/rgutils /rgcore/rgutils

COPY ./referee/go.mod  ./referee/go.sum /referee/

WORKDIR /referee

RUN go mod download
RUN go mod verify

COPY ./referee/internal internal
COPY ./referee/cmd cmd

# avoid dynamic linking for cgo
# -linkmode external
#   generate a standalone executable for cgo libraries 
#   instead of needing to refer to a dynamic library
# -extldflags "-static"
#   "-static" prevents linking with the shared libraries
# reduce binary size
# -s 
#   Omit the symbol table and debug information.
# -w 
#   Omit the DWARF symbol table.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-linkmode external -extldflags "-static" -w -s' -a -o cmd/main cmd/main.go

#-------------------------------------------------------
FROM scratch

COPY --from=backend-builder ./referee/cmd/main ./

EXPOSE 3333

ENTRYPOINT ["./main"]