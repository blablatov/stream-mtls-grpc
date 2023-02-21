module github.com/blablatov/stream-mtls-grpc

go 1.18

replace github.com/blablatov/stream-mtls-grpc/mtls-proto => ./mtls-proto

replace github.com/blablatov/stream-mtls-grpc/mockups => ./mockups

require (
	github.com/blablatov/stream-tls-grpc v0.0.0-20230219170727-7d2f5c5bfbff
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang/mock v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	golang.org/x/oauth2 v0.5.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)

require (
	cloud.google.com/go/compute v1.15.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)
