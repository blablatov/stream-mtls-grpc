[![Go](https://github.com/blablatov/stream-mtls-grpc/actions/workflows/stream-mtls-grpc.yml/badge.svg)](https://github.com/blablatov/stream-mtls-grpc/actions/workflows/stream-mtls-grpc.yml)
### Building and Running gRPC service  
Используется модель ошибок, встроенная в протокол gRPC и более развитая модель ошибок, реализованная в пакете `Google API google.rpc`.  
In order to build, Go to ``Go`` module directory location `stream-mtls-grpc/mtls-service` and execute the following shell command:  
```
go build -v 
./mtls-service
```   

### Building and Running gRPC client     
Используется модель ошибок, встроенная в протокол gRPC и более развитая модель ошибок, реализованная в пакете `Google API google.rpc`.  
In order to build, Go to ``Go`` module directory location `stream-mtls-grpc/mtls-client` and execute the following shell command:    
```
go build -v 
./mtls-client
```  

### Generates Server and Client side code via proto-file  
Go to ``Go`` module directory location `stream-mtls-grpc/mtls-proto` and execute the following shell commands:    
``` 
protoc product_info.proto --go_out=./ --go-grpc_out=./
protoc product_info.proto --go-grpc_out=require_unimplemented_servers=false:.
``` 
