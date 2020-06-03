# Kubernetes

## Config

kubectl create namespace go-user-microservice

kubectl config set-context --current --namespace go-user-microservice

kubectl create configmap api-server-config --from-env-file=.env
kubectl create configmap api-gw-config --from-env-file=.env
