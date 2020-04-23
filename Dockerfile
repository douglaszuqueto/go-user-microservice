FROM golang:alpine as builder
ARG service
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags "${XFLAGS} -s -w" -a -o service ./cmd/$service/$service.go

FROM alpine as upx
RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/service /app
RUN upx /app/service

FROM alpine
RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
COPY certs /app/certs
WORKDIR /app
COPY --from=upx /app/service /app
CMD ["./service"]
