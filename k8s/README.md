# Kubernetes

## Config

kubectl create namespace go-user-microservice

kubectl config set-context --current --namespace go-user-microservice

kubectl create configmap api-server-config --from-env-file=./server/.env
kubectl create configmap api-gw-config --from-env-file=./gw/.env

## Testes - SIEGE

Contexto real: PostgreSQL + GRPC + GRPC Gateway

Réplicas:

* GRPC = 1
* GW = 1

Obs: Sem limite nos resources - CPU/Memory

```bash
siege -t3600s -c50 http://k8s-server:30181/v1/user -v

Transactions:		     2281275 hits
Availability:		      100.00 %
Elapsed time:		     2297.80 secs
Data transferred:	      372.03 MB
Response time:		        0.05 secs
Transaction rate:	      992.81 trans/sec
Throughput:		        0.16 MB/sec
Concurrency:		       47.07
Successful transactions:     2281275
Failed transactions:	           0
Longest transaction:	        0.81
Shortest transaction:	        0.00
```