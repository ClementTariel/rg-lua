#-------------------------------------------------------
FROM golang:alpine AS backend-builder
# Add gcc for cgo
RUN apk add build-base

COPY ./bouncer/go.mod  ./bouncer/go.sum /bouncer/
COPY ./rgcore/rgentities /rgcore/rgentities

WORKDIR /bouncer

RUN go mod download
RUN go mod verify

COPY ./bouncer/internal internal
COPY ./bouncer/cmd cmd

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
FROM node:lts-alpine as frontend-builder

COPY ./bouncer/internal/presentation /bouncer/internal/presentation
WORKDIR /bouncer/internal/presentation
RUN npm install
RUN npm run build

#-------------------------------------------------------
FROM scratch

COPY --from=backend-builder ./bouncer/cmd/main ./

COPY --from=frontend-builder /bouncer/internal/presentation/dist ./dist

EXPOSE 5555

ENTRYPOINT ["./main"]