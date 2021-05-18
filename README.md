go build -o ./dist \
eval $(minikube docker-env) \
docker build -t coiner/trader-provisioner:0.0.1 . 

minikube service --url postgres-postgresql \
minikube service --url trader-provisioner
