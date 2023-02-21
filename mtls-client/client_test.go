// This conventional test.
// Before his execute run grpc-server ./mtls-service/mtls-service

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"testing"
	"time"

	pb "github.com/blablatov/stream-mtls-grpc/mtls-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

// Conventional test that starts a gRPC client test the service with RPC.
// Традиционный тест, который запускает клиент для проверки удаленного метода сервиса.
func TestServer_AddProduct(t *testing.T) {
	tokau := oauth.NewOauthAccess(fetchToken())

	// Load the client certificates from disk
	// Создаем пары ключей X.509 непосредственно из ключа и сертификата сервера
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	// Генерируем пул сертификатов в нашем локальном удостоверяющем центре
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	// Добавляем клиентские сертификаты из локального удостоверяющего центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs")
	}

	// Указываем аутентификационные данные для транспортного протокола с помощью DialOption.
	opts := []grpc.DialOption{
		// Указываем один и тот же токен OAuth в параметрах всех вызовов в рамках одного соединения.
		// Если нужно указывать токен для каждого вызова отдельно, используем CallOption.
		grpc.WithPerRPCCredentials(tokau),
		// Указываем транспортные аутентификационные данные в виде параметров соединения
		// Поле ServerName должно быть равно значению Common Name, указанному в сертификате
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	// Contact the server and print out its response.
	name := "Sumsung S999"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2029"
	price := float32(777.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Calls remote method of AddProduct
	// Вызываем удаленный метод AddProduct
	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
	if err != nil { // Checks response. Проверяем ответ
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Res %s", r.Value)
}

// Тестирование производительности в цикле за указанное колличество итераций
func BenchmarkServer_AddProduct(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 25; i++ {
		tokau := oauth.NewOauthAccess(fetchToken())

		// Load the client certificates from disk
		// Создаем пары ключей X.509 непосредственно из ключа и сертификата сервера
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Fatalf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		// Генерируем пул сертификатов в нашем локальном удостоверяющем центре
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		// Добавляем клиентские сертификаты из локального удостоверяющего центра в сгенерированный пул
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append ca certs")
		}

		opts := []grpc.DialOption{
			grpc.WithPerRPCCredentials(tokau),
			// Указываем транспортные аутентификационные данные в виде параметров соединения
			// Поле ServerName должно быть равно значению Common Name, указанному в сертификате
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				ServerName:   hostname, // NOTE: this is required!
				Certificates: []tls.Certificate{certificate},
				RootCAs:      certPool,
			})),
		}

		conn, err := grpc.Dial(address, opts...) // Подключаемся к серверному приложению
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewProductInfoClient(conn)

		// Contact the server and print out its response.
		name := "Sumsung S999"
		description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2029"
		price := float32(777.0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Calls remote method of AddProduct
		// Вызываем удаленный метод AddProduct
		r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
		if err != nil { // Checks response. Проверяем ответ
			log.Fatalf("Could not add product: %v", err)
		}
		log.Printf("Res %s", r.Value)
	}
}
