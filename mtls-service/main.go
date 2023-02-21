package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"path/filepath"
	"strings"

	pb "github.com/blablatov/stream-mtls-grpc/mtls-proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	crtFile            = filepath.Join("..", "mcerts", "server.crt")
	keyFile            = filepath.Join("..", "mcerts", "server.key")
	caFile             = filepath.Join("..", "mcerts", "ca.crt")
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

const (
	port = ":50051"
)

func main() {
	log.SetPrefix("Server event: ")
	log.SetFlags(log.Lshortfile)

	// Считываем и анализируем открытый/закрытый ключи и создаем сертификат, чтобы включить TLS
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	// Генерируем пул сертификатов в удостоверяющем центре.
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the client certificates from the CA
	// Добавляем клиентские сертификаты из удостоверяющего центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append client certs")
	}

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections. Включаем TLS для всех входящих соединений путем.
		grpc.Creds( // Create the TLS credentials. Создание аутентификационных данных TLS.
			credentials.NewTLS(&tls.Config{
				ClientAuth:   tls.RequireAndVerifyClientCert,
				Certificates: []tls.Certificate{cert},
				ClientCAs:    certPool,
			},
			)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// Registers unary interceptor to gRPC-server
			// Будет направлять все клиентские запросы к функции ensureValidBasicCredentials
			grpc.UnaryServerInterceptor(ensureValidToken),
			// Регистрация дополнительного унарного перехватчика на gRPC-сервере
			// Будет направлять все клиентские запросы к функции orderUnaryServerInterceptor
			grpc.UnaryServerInterceptor(orderUnaryServerInterceptor),
		)),
	}

	// Создаем новый экземпляр gRPC-сервера, передавая ему аутентификационные данные
	s := grpc.NewServer(opts...)
	// Регистрируем реализованный сервис на только что созданном gRPCсервере с помощью сгенерированных AP
	pb.RegisterProductInfoServer(s, &server{})

	lis, err := net.Listen("tcp", port) // Начинаем прослушивать TCP-соединение на порту 50051.
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting gRPC listener on port " + port)

	// Привязываем gRPC-сервер к прослушивателю и ждем появления сообщений на порту 50051.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Performs validation of token matching an arbitrary string.
	// Выполняем проверку токена, соответствующего произвольной строке.
	return token == "blablatok-tokblabla-blablatok"
}

// Checking token. Определяем функцию ensureValidToken для проверки подлинности токена.
// Если тот отсутствует или недействителен, тогда перехватчик блокирует выполнение и возвращает ошибку.
// Или вызывается следующий обработчик, которому передается контекст и интерфейс.
func ensureValidToken(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

// Server : Unary Interceptor
// Серверный унарный перехватчик в gRPC
func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pre-processing logic
	// Gets info about the current RPC call by examining the args passed in
	// Логика перед вызовом. Получает информацию о текущем RPC-вызове путем анализа переданных аргументов
	log.Println("====== [Server Interceptor] ", info.FullMethod)
	log.Printf(" Pre Proc Message : %s", req)

	// Invoking the handler to complete the normal execution of a unary RPC.
	// Вызываем обработчик, чтобы завершить нормальное выполнение унарного RPC-вызова
	m, err := handler(ctx, req)

	// Post processing logic
	// Логика после вызова
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}
