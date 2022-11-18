## Build
FROM golang:1.17-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /KineJameAPI

## Deploy
FROM alpine:latest

WORKDIR /

COPY --from=build /KineJameAPI /KineJameAPI
ENTRYPOINT [ "/KineJameAPI" ]
CMD [ "/KineJameAPI" ] 