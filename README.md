# filebeat output loki

`filebeat` official release not support output `loki` via grafana

## use
```shell
docker pull ghcr.io/wolanx/filebeat:main
```

## config demo

```yml
# http
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - C:\Users\admin\Desktop\pic\*.log

output.loki:
  hosts: [ 'localhost:3100' ]
  is_grpc: false
```

```yml
# grpc recommend
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - C:\Users\admin\Desktop\pic\*.log

output.loki:
  hosts: [ 'localhost:9095' ]
  is_grpc: true
```
