compile:
	if [ -e .dist/trader ];	then rm dist/trader; fi;
	GOOS=linux go build -o dist/trader .

remove-build-docker:
	eval $(minikube docker-env)
	if [[ "$(docker images -q coiner/trader-provisioner:0.0.1 2> /dev/null)" != "" ]]; then docker image rm coiner/trader-provisioner:0.0.1; fi;
	docker build -t coiner/trader-provisioner:0.0.1 .

dockerize: compile remove-build-docker
