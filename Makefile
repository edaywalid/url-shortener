runm :
	docker-compose up -d --build

stopm :
	docker-compose down

run :
	docker run --name my-zookeeper -p 2181:2181 -d zookeeper
	docker run --name my-redis -p 6379:6379 -d redis
	go run cmd/server/main.go

stop :
	docker stop my-zookeeper
	docker stop my-redis
	docker rm my-zookeeper
	docker rm my-redis
