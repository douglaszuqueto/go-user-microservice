#!/bin/bash

docker build -t douglaszuqueto/go-user-microservice-$1 --build-arg service=$1 .

docker push douglaszuqueto/go-user-microservice-$1