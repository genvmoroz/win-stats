cd picker
go fmt ./...
go get -u -t ./...
go mod tidy
go build -o ./app.exe ./cmd/service
start /MIN app.exe ^& exit
exit
