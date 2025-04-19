cd prometheus-collector
go fmt ./...
go get -u -t ./...
go mod tidy
cd deployment
docker-compose rm -f
docker-compose pull
docker-compose up --force-recreate --build -d
exit
