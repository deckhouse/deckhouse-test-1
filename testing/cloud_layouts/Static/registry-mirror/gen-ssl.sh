#! /bin/bash

rm -rf ./ssl && mkdir ssl
openssl req -x509 -subj=/C=RU/ST=Moscow/L=Moscow/O=flant/OU=IT -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.pem
mv server* ssl