# Multi-stage tls-server build
# Многоэтапная сборка tls-service

FROM golang AS build

ENV location /go/src/github.com/blablatov/stream-tls-grpc

WORKDIR ${location}/tls-service

ADD ./main.go ${location}/tls-service
ADD ./tls-proto  ${location}/tls-proto

#ADD ./go.mod  ${location}/tls-service
#ADD ./go.sum  ${location}/tls-service
RUN go mod init github.com/blablatov/stream-tls-grpc/tls-service

# Download all the dependencies. Загрузка зависимостей
RUN go get -d ./...
# Install the package. Установка пакетов
RUN go install ./...

# Build server app. Сборка серверного приложения.
RUN CGO_ENABLED=0 go build -o /bin/tls-service 

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.
FROM scratch
COPY --from=build /bin/tls-service /bin/tls-service

ENTRYPOINT ["/bin/tls-service"]
EXPOSE 50051