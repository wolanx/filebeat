# filebeat output loki

`filebeat` official built not support output `loki` via grafana

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
```

```yml
# grpc
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - C:\Users\admin\Desktop\pic\*.log

output.loki:
  hosts: [ 'localhost:9095' ]
```
