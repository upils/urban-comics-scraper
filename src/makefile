NAME=ucscraper

build:
	go build -o bin/${NAME}

build_docker_linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${NAME}

build_docker_macos:
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o ${NAME}

run:
	./bin/${NAME}

clean:
	rm ./bin/${NAME}
