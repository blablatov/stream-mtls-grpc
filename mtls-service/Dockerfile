# Multi-stage mtls-server build
# Многоэтапная сборка mtls-service

FROM golang AS build

ENV location /go/src/github.com/blablatov/stream-mtls-grpc

WORKDIR ${location}/mtls-service

ADD ./main.go ${location}/mtls-service
ADD ./mtls-proto  ${location}/mtls-proto

#ADD ./go.mod  ${location}/mtls-service
#ADD ./go.sum  ${location}/mtls-service
RUN go mod init github.com/blablatov/stream-mtls-grpc/mtls-service

# Download all the dependencies. Загрузка зависимостей
RUN go get -d ./...
# Install the package. Установка пакетов
RUN go install ./...

# Build server app. Сборка серверного приложения.
RUN CGO_ENABLED=0 go build -o /bin/mtls-service 

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.
FROM scratch
COPY --from=build /bin/mtls-service /bin/mtls-service

ENTRYPOINT ["/bin/mtls-service"]
EXPOSE 50051