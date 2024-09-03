cd prometheus-collector
go get -u -t ./...
make all
cd deployment
docker-compose rm -f
docker-compose pull
docker-compose up --force-recreate --build --abort-on-container-exit
