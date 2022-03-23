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