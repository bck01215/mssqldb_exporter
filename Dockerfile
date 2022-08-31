FROM golang:1.18-alpine3.15 as build

RUN mkdir -p /app
COPY go.mod go.sum /app/
COPY sql-server-exporter/*.go /app/sql-server-exporter/
WORKDIR /app/sql-server-exporter
RUN go build -o sql-server-exporter
RUN chmod a+x ./sql-server-exporter

FROM alpine
COPY --from=build /app/sql-server-exporter/sql-server-exporter /app/sql-server-exporter
WORKDIR /app/
EXPOSE 9101
ENTRYPOINT [ "/app/sql-server-exporter"]