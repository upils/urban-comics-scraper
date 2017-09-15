NAME=ucscraper

build:
	go build -o bin/${NAME}

build_docker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${NAME}

all: build run

run:
	./bin/${NAME}

clean:
	rm ./bin/${NAME}
