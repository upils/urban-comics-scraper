#!/bin/bash

echo -e "\e[32m[+] Setup the docker go builder.\e[0m"

docker-compose build

echo -e "\e[32m[+] Compile the binary.\e[0m"

docker-compose run

echo -e "\e[32m[+] Clean containers.\e[0m"

docker-compose rm -v
