FROM golang:1.16
RUN mkdir -p /app
COPY . /app
VOLUME ["/app"]
WORKDIR /app/sql-server-exporter
RUN go build && mv sql-server-exporter /
RUN ls
WORKDIR /app
RUN rm -rf *
EXPOSE 9101
ENTRYPOINT ["/sql-server-exporter"]