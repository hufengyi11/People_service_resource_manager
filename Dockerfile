# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .

RUN go get ./...

RUN go build -o /people_service

FROM alpine
COPY --from=build /people_service .
EXPOSE 8080
ENTRYPOINT [ "./people_service" ]