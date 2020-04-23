#!/bin/bash

docker build -t grpc-$1 --build-arg service=$1 .