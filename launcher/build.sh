CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/imserver main.go


cp -rf scripts/* ./build/
