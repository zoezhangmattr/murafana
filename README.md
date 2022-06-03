# murafana
a small app to query grafana api

## get all dashboards
get all dashboard and save to yaml file

## usage
* add .env file to local, source .env
```.env
export GRAFANA_CLOUD_API_KEY="xxxxx"
export GRAFANA_URL="https://xxxx"
```
* run command
```sh
go run main.go
go run main.go -uid AWSWAFV2a
go run main.go -download-list enabled
```
