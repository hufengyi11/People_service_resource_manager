# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /people-service


CMD [ "/people-service" ]

# # RUN go get ./...

# RUN go build -o /people-service

# ## Deploy
# FROM gcr.io/distroless/base-debian10

# WORKDIR /

# COPY --from=build /docker-gs-ping /docker-gs-ping

# EXPOSE 8080

# USER nonroot:nonroot

# ENTRYPOINT ["/people-service"]