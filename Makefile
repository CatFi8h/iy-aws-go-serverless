build:
	echo "Building lambda binaries"
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/create/bootstrap functions/create/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/update/bootstrap functions/update/main.go
	env GOOS=linux GOARCH=arm64 go build -o build/lambda/get/bootstrap functions/get/main.go

zip:
	zip -j build/lambda/create.zip build/lambda/create/bootstrap
	zip -j build/lambda/update.zip build/lambda/update/bootstrap
	zip -j build/lambda/get.zip build/lambda/get/bootstrap

clear:
	rm -rf build

deploy:
	serverless deploy --verbose