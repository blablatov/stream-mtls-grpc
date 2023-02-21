# Multi-stage mtls-client build
# Многоэтапная сборка mtls-client

FROM golang AS build

ENV location /go/src/github.com/blablatov/stream-mtls-grpc

WORKDIR ${location}/mtls-client

ADD ./main.go ${location}/mtls-client

RUN go mod init github.com/blablatov/stream-mtls-grpc/mtls-client

RUN CGO_ENABLED=0 go build -o mtls-client

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.
FROM scratch
COPY --from=build ./mtls-client ./mtls-client

ENTRYPOINT ["./mtls-client"]
EXPOSE 50051