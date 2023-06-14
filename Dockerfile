FROM golang:1.20 AS build

USER root
WORKDIR /app
COPY main.go main.go
COPY go.mod  go.mod

RUN go get && go mod tidy && go build main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o metrics main.go

FROM alpine:3.14 AS run

WORKDIR /app
COPY --from=build /app/metrics /app/metrics

EXPOSE 80
ENTRYPOINT [ "/app/metrics" ]