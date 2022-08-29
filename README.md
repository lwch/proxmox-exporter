# proxmox-exporter

proxmox exporter of prometheus, you must install the exporter in each proxmox node.

* supported node info
* supported vm info
* supported sensors info
* supported disk info

view [metrics](docs/metrics.txt)

## how to install

1. download exporter file from [latest](https://github.com/lwch/proxmox-exporter/releases/latest) version, and add execute permission

       sudo chmod +x exporter
2. create configure file `exporter.yaml` by [example](https://github.com/lwch/proxmox-exporter/blob/master/conf/exporter.yaml)
3. (optional)create token from web => datacenter => permissions => api tokens
4. change api.user and api.token from configure file
5. use command to install linux service

       sudo ./exporter -conf exporter.yaml -action install
6. use systemctl command to start service

       sudo systemctl start proxmox-exporter

7. add node in prometheus
8. import dashboard in grafana

## grafana dashboard

![grafana](docs/grafana.png)

## supported

only tested proxmox version 7.2