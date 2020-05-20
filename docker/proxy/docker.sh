#!/bin/bash

docker run -it --rm \
    -p 8080:8080 \
    -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro \
    -v $(pwd)/vh/default.conf:/etc/nginx/conf.d/default.conf:ro \
    nginx:alpine