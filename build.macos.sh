#!/bin/bash

echo -e "\033[32m[+] Run builder\033[0m"
docker-compose run builder make build_docker_macos
echo -e "\033[32m[+] Clean containers\033[0m"
docker-compose down