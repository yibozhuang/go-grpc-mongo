.DEFAULT_GOAL := default

MONGO_CHECK := $(shell lsof -n -i4TCP:27017 | grep LISTEN)

prep:
	mkdir -p ./build

build/server:
	go build -o build/server server/main.go

build/blog:
	go build -o build/blog ./client

mongo:
ifeq ("$(MONGO_CHECK)", "")
	docker-compose up -d
endif

stop_mongo:
	docker-compose down

run_server:
	./build/server &

stop_server:
	killall server

run_client:
	./build/blog create -a 1 -c "This is test1" -t "Test1"
	./build/blog list

test: mongo run_server run_client stop_server

default: prep build/server build/blog

clean: stop_mongo
	rm -rf ./build
