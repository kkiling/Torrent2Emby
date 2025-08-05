# buf build

### install
Установка buf build
```
curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m).tar.gz" | \
sudo tar -xvzf - -C /usr/local --strip-components 1

buf --version
```
Установка
```
sudo apt install protobuf-compiler

protoc --version  
```

Добавляем в go.mod
```
go 1.24

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)

go mod tidy
go install tool
```

Generate
```
buf dep update
buf generate
```

Если buf 403
```
export HTTP_PROXY="http://USER:PASSWORD@IP:PORT"
export HTTPS_PROXY="http://USER:PASSWORD@IP:PORT"
proxychains buf generate
```