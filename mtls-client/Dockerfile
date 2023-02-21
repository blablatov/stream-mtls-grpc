# Multi-stage tls-client build
# Многоэтапная сборка tls-client

FROM golang AS build

ENV location /go/src/github.com/blablatov/stream-tls-grpc

WORKDIR ${location}/tls-client

ADD ./main.go ${location}/tls-client

RUN go mod init github.com/blablatov/stream-tls-grpc/tls-client

RUN CGO_ENABLED=0 go build -o tls-client

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.
FROM scratch
COPY --from=build ./tls-client ./tls-client

ENTRYPOINT ["./tls-client"]
EXPOSE 50051