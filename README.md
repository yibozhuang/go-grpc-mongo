# go-grpc-mongo
Experiment with Go gRPC with Mongo

### Generating package from protobuf
```
$ protoc proto/blog.proto --go_out=plugins=grpc:.
```

### Starting local mongo instance
```
$ docker-compose up -d
```

### Running server
```
$ go run server/main.go
```

### Compiling and running server
```
$ go build -o build/server server/main.go
$ ./build/server
```

### Running client
```
$ go run client/main.go
```

### Compiling and running client
```
$ go build -o build/blog ./client
$ ./build/blog create -a 1 -c "This is test1" -t "Test1"
```

### Running everything
```
$ make clean
$ make
$ make test
```
