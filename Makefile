build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

docker-build-run:
	echo "building docker file"
	@docker build -t go-hotel-api .
	echo "running API inside Docker container"
	@docker run -p 3000:3000 go-hotel-api

docker-run:
	echo "running API inside Docker container"
	@docker run -p 5000:5000 go-hotel-api

test:
	@go test -v ./...