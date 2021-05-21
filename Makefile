build:
	rm dist/trader
	go build -o dist/trader .

docker:
	docker image rm coiner/trader-provisioner:latest
	docker build -t coiner/trader-provisioner:latest .