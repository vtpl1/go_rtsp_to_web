```bash
cd vtpl1/vdk/format/vtplgrpc 
```
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative stream_service.proto
```
```go
go get google.golang.org/grpc
go get google.golang.org/grpc/codes
go get google.golang.org/grpc/status
```