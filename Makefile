compile:
	rm dist/trader
	GOOS=linux go build -o dist/trader .

remove-build-docker:
	docker image rm coiner/trader-provisioner:latest
	docker build -t coiner/trader-provisioner:latest .

dockerize: compile remove-build-docker
