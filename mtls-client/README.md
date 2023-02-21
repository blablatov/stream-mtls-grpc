### Тестирование функциональность клиентского кода с подключением к серверу. 
### Testing code with conn to server          
  
Традиционный тест, который запускает клиент для проверки удаленного метода сервиса.  
Перед его выполнением запустить grpc-сервер.   
(Conventional test that starts a gRPC client test the service with RPC.Before his execute run grpc-server):      

```shell script
./tls-service/tls-service
```

Для тестирования клиента, без подключения к серверу, выполнить сгенерированный тестовый код.      
(Runs generation code of mock up for interface ProductInfoClient):   
       
```shell script
./mockups/prodinfo_mock_test.go
```

### Create and run Docker image

Создание Docker контейнера для gRPC-клиента (build container of client):    

```shell script
docker build -t tls-client .
```

Развернуть задание с клиентским gRPC-приложением:    

```shell script
kubectl apply -f grpc-tls-client.yaml
```


