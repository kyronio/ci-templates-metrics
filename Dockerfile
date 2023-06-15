FROM golang:1.20 AS build

USER root
WORKDIR /app
COPY . .

RUN go get && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o metricsparser main.go

FROM alpine:3.14 AS run

ENV RECEIVER_SERVICE=http://metrics:80
ENV METRIC_PATH=/build/.metrics/metrics.json

WORKDIR /app
COPY --from=build /app/metricsparser /app/metricsparser

EXPOSE 80
ENTRYPOINT [ "/app/metricsparser" ]