FROM alpine:latest

WORKDIR /app
ADD etl /usr/local/bin

CMD ["etl", "server"]
