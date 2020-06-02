#  BASE
FROM golang:alpine as base
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

# BUILDER
FROM base as builder
ARG service
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "${XFLAGS} -s -w" -a -o service ./cmd/$service/$service.go

# UPX
FROM douglaszuqueto/alpine-upx as upx
WORKDIR /app
COPY --from=builder /app/service /app
RUN upx /app/service

# FINAL
FROM douglaszuqueto/alpine-go
COPY certs /app/certs
WORKDIR /app
COPY --from=upx /app/service /app
CMD ["./service"]
