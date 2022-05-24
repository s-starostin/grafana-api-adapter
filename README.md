# Grafana API adapter

You could specify the absolute path to the directory with config.ini via the GRAFANA_ADAPTER_CONF environment variable:
```bash
export GRAFANA_ADAPTER_CONF=/tmp && ./grafana-adapter & 
```
## config.ini
```
[server]
PORT = 8000
HOST = 0.0.0.0
LOGIN = grafana-adapter
PASSWORD = s0m3passw0rD
[grafana]
PORT = 3000
HOST = localhost
LOGIN = admin
PASSWORD = Somepa$$w0rd
```

Block `server` specifies grafana-adapter server parameters: 

| Parameter | Type | Default | Comment |
| ------ | ------ | -------  | ---------- |
| PORT | int | 80 | Server port |
| HOST | string | localhost | Server host |
| LOGIN | string | "" | Basic auth login |
| PASSWORD | string | "" | Basic auth password |

Block `grafana` specifies grafana instance parameters to connect to: 

| Parameter | Type | Default | Comment |
| ------ | ------ | -------  | ---------- |
| PORT | int | 3000 | Grafana port |
| HOST | string | localhost | Grafana host |
| LOGIN | string | admin | Grafana admin login |
| PASSWORD | string | admin | Grafana admin password |

## API
### Users
Retrieving all:
```
GET
.../users
```

Retrieving | deleting single user:
```
GET | DELETE
.../users/{id*} (.../users/1 || .../users/id=1 || .../users/login=admin || .../users/email=admin@localhost)
.../users/?id=1
.../users/?login=test
.../users/?email=test@test.test
```

Creating user:
```
POST
.../users/ (data: {})
```

examples:
```
curl -X POST adapter:8000/users/ -H 'Content-Type: application/json' -d '{"email":"test@test.test"}'
curl -X POST adapter:8000/users/ -H 'Content-Type: application/json' -d '{"email":"test@test.test", "login":"test", "password":"t3st"}'
```

Adding user to specific organizations:
```
PATCH
.../users/organizations/ (data: {})
```

examples:
```
curl -X POST adapter:8000/users/organizations/ -H 'Content-Type: application/json' -d '{"user":{"login":"test","email":"test@test.test"},"organizations":[{"name":"test","role":"Admin"},{"name":"test2","role":"Editor"}]}'
curl -X POST adapter:8000/users/organizations/ -H 'Content-Type: application/json' -d '{"user":{"email":"test@test.test"},"organizations":[{"name":"test","role":"Admin"}]}'
```

### Organizations
Retrieving all:
```
GET
.../organizations
```

Retrieving | deleting single organization:
```
GET | DELETE
.../organizations/{id*} (.../organizations/1 || .../organizations/test)
.../organizations/?id=1
.../organizations/?name=test
```

Creating organization:
```
POST
.../organizations/ (data: {})
```

examples:
```
curl -X POST adapter:8000/organizations/ -H 'Content-Type: application/json' -d '{"name":"test"}'
```

### Dashboards for organization
Retrieving all:
```
GET
.../organizations/{orgId}/dashboards/ (.../organizations/11/dashboards/ || .../organizations/test/dashboards/)
```

Retrieving | deleting single dashboard:
```
GET | DELETE
.../organizations/{orgId}/dashboards/{uid} (.../organizations/11/dashboards/GPXicXZRk || .../organizations/test/dashboards/organization%20title || .../organizations/test/dashboards/23)
```

Creating | updating dashboard:
```
POST
.../organizations/{orgId}/dashboards/ (data: {})
```

creating examples:
```
curl -X POST adapter:8000/organizations/1/dashboards/ -H 'Content-Type: application/json' -d '{"dashboard":{"title":"test"}}'
curl -X POST adapter:8000/organizations/1/dashboards/ -H 'Content-Type: application/json' -d '{"dashboard":{"annotations":{"list":[{"builtIn":1,"datasource":"-- Grafana --","enable":true,"hide":true,"iconColor":"rgba(0, 211, 255, 1)","name":"Annotations & Alerts","target":{"limit":100,"matchAny":false,"tags":[],"type":"dashboard"},"type":"dashboard"}]},"editable":true,"gnetId":null,"graphTooltip":0,"links":[],"panels":[{"datasource":null,"fieldConfig":{"defaults":{"color":{"mode":"palette-classic"},"custom":{"axisLabel":"","axisPlacement":"auto","barAlignment":0,"drawStyle":"line","fillOpacity":0,"gradientMode":"none","hideFrom":{"legend":false,"tooltip":false,"viz":false},"lineInterpolation":"linear","lineWidth":1,"pointSize":5,"scaleDistribution":{"type":"linear"},"showPoints":"auto","spanNulls":false,"stacking":{"group":"A","mode":"none"},"thresholdsStyle":{"mode":"off"}},"mappings":[],"thresholds":{"mode":"absolute","steps":[{"color":"green","value":null},{"color":"red","value":80}]}},"overrides":[]},"gridPos":{"h":8,"w":12,"x":0,"y":0},"id":2,"options":{"legend":{"calcs":[],"displayMode":"list","placement":"bottom"},"tooltip":{"mode":"single"}},"title":"Panel Title","type":"timeseries"}],"schemaVersion":30,"style":"dark","tags":[],"templating":{"list":[]},"time":{"from":"now-6h","to":"now"},"timepicker":{},"timezone":"","title":"test2"}}'
```

updating examples (id/uid):
```
curl -X POST adapter:8000/organizations/1/dashboards/ -H 'Content-Type: application/json' -d '{"dashboard":{"annotations":{"list":[{"builtIn":1,"datasource":"-- Grafana --","enable":true,"hide":true,"iconColor":"rgba(0, 211, 255, 1)","name":"Annotations & Alerts","target":{"limit":100,"matchAny":false,"tags":[],"type":"dashboard"},"type":"dashboard"}]},"editable":true,"gnetId":null,"graphTooltip":0,"links":[],"panels":[{"datasource":null,"fieldConfig":{"defaults":{"color":{"mode":"palette-classic"},"custom":{"axisLabel":"","axisPlacement":"auto","barAlignment":0,"drawStyle":"line","fillOpacity":0,"gradientMode":"none","hideFrom":{"legend":false,"tooltip":false,"viz":false},"lineInterpolation":"linear","lineWidth":1,"pointSize":5,"scaleDistribution":{"type":"linear"},"showPoints":"auto","spanNulls":false,"stacking":{"group":"A","mode":"none"},"thresholdsStyle":{"mode":"off"}},"mappings":[],"thresholds":{"mode":"absolute","steps":[{"color":"green","value":null},{"color":"red","value":80}]}},"overrides":[]},"gridPos":{"h":8,"w":12,"x":0,"y":0},"id":2,"options":{"legend":{"calcs":[],"displayMode":"list","placement":"bottom"},"tooltip":{"mode":"single"}},"title":"Panel Title","type":"timeseries"}],"schemaVersion":30,"style":"dark","tags":[],"templating":{"list":[]},"time":{"from":"now-6h","to":"now"},"timepicker":{},"timezone":"","title":"test2","id":29}}'
curl -X POST adapter:8000/organizations/1/dashboards/ -H 'Content-Type: application/json' -d '{"dashboard":{"annotations":{"list":[{"builtIn":1,"datasource":"-- Grafana --","enable":true,"hide":true,"iconColor":"rgba(0, 211, 255, 1)","name":"Annotations & Alerts","target":{"limit":100,"matchAny":false,"tags":[],"type":"dashboard"},"type":"dashboard"}]},"editable":true,"gnetId":null,"graphTooltip":0,"links":[],"panels":[{"datasource":null,"fieldConfig":{"defaults":{"color":{"mode":"palette-classic"},"custom":{"axisLabel":"","axisPlacement":"auto","barAlignment":0,"drawStyle":"line","fillOpacity":0,"gradientMode":"none","hideFrom":{"legend":false,"tooltip":false,"viz":false},"lineInterpolation":"linear","lineWidth":1,"pointSize":5,"scaleDistribution":{"type":"linear"},"showPoints":"auto","spanNulls":false,"stacking":{"group":"A","mode":"none"},"thresholdsStyle":{"mode":"off"}},"mappings":[],"thresholds":{"mode":"absolute","steps":[{"color":"green","value":null},{"color":"red","value":80}]}},"overrides":[]},"gridPos":{"h":8,"w":12,"x":0,"y":0},"id":2,"options":{"legend":{"calcs":[],"displayMode":"list","placement":"bottom"},"tooltip":{"mode":"single"}},"title":"Panel Title","type":"timeseries"}],"schemaVersion":30,"style":"dark","tags":[],"templating":{"list":[]},"time":{"from":"now-6h","to":"now"},"timepicker":{},"timezone":"","title":"test22","uid": "ITm_ajWgk"}}'
```

### Folders for organization
Retrieving all:
```
GET
.../organizations/{orgId}/folders (.../organizations/11/dashboards/)
```

Retrieving | deleting single folder:
```
GET | DELETE
.../organizations/{orgId}/folders/{id} (.../organizations/11/folders/nErXDvCkzz || .../organizations/11/folders/11)
```

Creating folder:
```
POST
.../organizations/{orgId}/folders/ (data: {})
```

creating examples:
```
curl -X POST adapter:8000/organizations/1/folders/ -H 'Content-Type: application/json' -d '{"title":"test"}'
```

### Datasources for organization
Retrieving all:
```
GET
.../organizations/{orgId}/datasources (.../organizations/11/datasources/)
```

Retrieving | deleting single datasource:
```
GET | DELETE
.../organizations/{orgId}/datasources/{id} (.../organizations/11/datasources/1 || .../organizations/test/datasources/test%20dashboard)
```

Creating datasource:
```
POST
.../organizations/{orgId}/datasources/ (data: {})
```

creating examples:
```
curl -X POST adapter:8000/organizations/test/datasources/ -H 'Content-Type: application/json' -d '{"name":"test","access":"proxy","type":"prometheus","jsonData":{"customQueryParameters":"test","httpMethod":"POST"}}'
```