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
  hosts: [ 'svc-loki:3100' ]
  protocol: http
```

```yml
# grpc recommend default
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - C:\Users\admin\Desktop\pic\*.log

output.loki:
  hosts: [ 'svc-loki:9095' ]
  protocol: grpc
```
