## Service and Client - Go Implementation

## Building and Running Service  

Клиентское и серверное приложения обмениваются сообщениями с префиксом длины, не дожидаясь завершения взаимодействия с противоположной стороны. 
Клиент и сервер отправляют сообщения одновременно.  
Любой из них может закрыть соединение на своей стороне, теряя тем самым возможность отправлять дальнейшие сообщения.  

In order to build, Go to ``Go`` module directory location `stream-mtls-grpc/mtls-service` and execute the following
 shell command:
```
go build -v 
./mtls-server
```  

## Building and Running Client   

In order to build, Go to ``Go`` module directory location `stream-mtls-grpc/mtls-client` and execute the following shell command:
```
go build -v 
./mtls-client
```

## Additional Information

### Generate Server and Client side code   
Go to ``Go`` module directory location `stream-mtls-grpc/mtls-proto` and execute the following shell commands:    
``` 
protoc order_management.proto --go_out=./ --go-grpc_out=./
protoc order_management.proto --go-grpc_out=require_unimplemented_servers=false:.
``` 
