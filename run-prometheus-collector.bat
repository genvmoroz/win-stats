REM cd prometheus-collector
REM go fmt ./...
REM go get -u -t ./...
REM go mod tidy
cd deployment
docker-compose rm -f
docker-compose pull
REM docker-compose up --force-recreate --build -d
docker-compose up -d
exit
