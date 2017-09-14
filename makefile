NAME=ucscraper

build:
	go build -o bin/${NAME}

docker_binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/${NAME}

all: build run

run:
	./bin/${NAME}

clean:
	rm ./bin/${NAME}
