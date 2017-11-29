#!/bin/bash

echo -e "\e[32m[+] Setup the docker go builder\e[0m"
docker-compose build
echo -e "\e[32m[+] Compile the binary\e[0m"
docker-compose up
echo -e "\e[32m[+] Clean containers\e[0m"
docker-compose down
echo -e "\e[32m[+] Changing own of the binary\e[0m"
sudo chown  $USER:$USER src/ucscraper

