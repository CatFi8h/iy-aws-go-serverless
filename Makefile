.PHONY: clear build zip deploy

build:
	echo "Building lambda binaries"
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/create/bootstrap cmd/functions/createDeviceInfo/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/update/bootstrap cmd/functions/updateDeviceInfo/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/delete/bootstrap cmd/functions/deleteDeviceInfo/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/get/bootstrap cmd/functions/getDeviceInfo/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/listensns/bootstrap cmd/functions/listenSns/main.go

zip:
	zip -j build/lambda/create.zip build/lambda/create/bootstrap
	zip -j build/lambda/update.zip build/lambda/update/bootstrap
	zip -j build/lambda/delete.zip build/lambda/delete/bootstrap
	zip -j build/lambda/get.zip build/lambda/get/bootstrap
	zip -j build/lambda/listensns.zip build/lambda/listensns/bootstrap

clear:
	rm -rf build

deploy:
	serverless deploy --verbose

testdynamodb:
	docker run -p 8000:8000 amazon/dynamodb-local

createTable:
	aws dynamodb create-table --table-name device-info-dev --attribute-definitions AttributeName=deviceId,AttributeType=S --key-schema AttributeName=deviceId,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://localhost:8000
	aws dynamodb create-table --attribute-definitions AttributeName=deviceid,AttributeType=S --table-name device-info-table --key-schema AttributeName=deviceid,KeyType=HASH --provisioned-throughput ReadCapacityUnits=2,WriteCapacityUnits=2 --endpoint-url http://localhost:32813

