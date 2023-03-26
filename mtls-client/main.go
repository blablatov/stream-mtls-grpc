package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	pb "github.com/blablatov/stream-mtls-grpc/mtls-proto"
	"golang.org/x/oauth2"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
)

var (
	crtFile = filepath.Join("..", "mcerts", "client.crt")
	keyFile = filepath.Join("..", "mcerts", "client.key")
	caFile  = filepath.Join("..", "mcerts", "ca.crt")
)

const (
	address = "localhost:50051"
	//address  = "net-tls-service:50051"
	hostname = "localhost"
)

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// Set up the credentials for the connection.
	// Значение токена OAuth2. Используем строку, прописанную в коде.
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

	// Set up a connection to the server.
	// Устанавливаем безопасное соединение с сервером, передавая параметры аутентификации
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Sends connect, it include all methods of server. Передаем соединение и создаем заглушку.
	// Ее экземпляр содержит все удаленные методы, которые можно вызвать на сервере.
	client := pb.NewProductInfoClient(conn)

	// Add invalid Order. Этот ID заказа недействителен
	// Contact the server and print out its response. Отправка данных на сервер, получение ответа.
	name := "-1"
	description := "Samsung Galaxy S000 not famous smartphone, launched now 0000"
	price := float32(0001.0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Calls remote method Add Order and assigns an error him
	// Вызываем удаленный метод AddOrder и присваиваем ошибку переменной addOrderError.
	res, addOrderError := client.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price}, grpc.UseCompressor(gzip.Name))

	if addOrderError != nil {
		errorCode := status.Code(addOrderError) // Get code an error. Получаем код ошибки из пакета grpc/status.
		if errorCode == codes.InvalidArgument { // Compares an error. Сравниваем код ошибки с InvalidArgument.
			log.Printf("Invalid Argument Error : %s", errorCode)
			errorStatus := status.Convert(addOrderError) // Gets of status. Получаем состояние.
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				// Checks type of error. Проверяем, имеет ли ошибка тип BadRequest_FieldViolation.
				case *epb.BadRequest_FieldViolation:
					log.Printf("Request Field Invalid: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error : %s ", errorCode)
		}
	} else {
		log.Print("AddOrder Response -> ", res.Value)
	}

	// Data for add. Contact the server and print out its response.
	// Отправка достоверных данных на сервер, получение ответа.
	name = "Sumsung S9999"
	description = "Samsung Galaxy S9999 is the latest smart phone, launched in February 2039"
	price = float32(7777.0)

	// Add Order. Добавляет заказ на сервере.
	r, err := client.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	// Response of server. Ответ сервера об успешном добавлении заказа с его номером.
	product, err := client.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Println("Product: ", product.String())
}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "blablatok-tokblabla-blablatok",
	}
}
