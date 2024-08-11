.PHONY: build deploy clean intial all all-swap

build:
	GOOS=linux GOARCH=amd64 go build -o ./pbr-service/bootstrap ./pbr-service/cmd/main.go  

deploy:
	cd deploy-scripts && cdk deploy

deploy-swap:
	cd deploy-scripts && cdk deploy --hotswap

clean:
	rm -rf ./pbr-service/bootstrap
	
intial:
	cd deploy-scripts && cdk bootstrap
	cd deploy-scripts && cdk synth
	make build
	make deploy

all:
	make clean
	make build
	make deploy

all-swap:
	make clean
	make build
	make deploy-swap