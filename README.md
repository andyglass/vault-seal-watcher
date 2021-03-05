# Vault Seal Watcher

## Usage

To start docker image run:
```bash
docker run --rm -p 8000:8000 krezz/vault-seal-watcher
```

To start docker-compose stack run:
```bash
docker-compose up --build
```


## Environment variables

| Environment variable | Default value | Type | Description |
| --- | --- | --- | --- |
| `SERVER_HOST` | `0.0.0.0` | `string` | Web server listen host |
| `SERVER_PORT` | `8000` | `int` | Web server listen port |
| `SERVER_LOG_LEVEL` | `info` | `string` | Log level |
| `SERVER_READ_TIMEOUT` | `10s` | `duration` | Web server read timeout |
| `SERVER_WRITE_TIMEOUT` | `10s` | `duration` | Web server write timeout |
| `VAULT_ADDR` | `http://localhost:8200` | `string` | Vault Http(s) address |
| `VAULT_TIMEOUT` | `10s` | `duration` | Timeout on connect to Vault server |
| `VAULT_UNSEAL_KEYS` | `nil` | `string` | Comma-separated list of key-shares |
| `VAULT_UNSEAL_DELAY` | `200ms` | `duration` | Delay between using key-shares to unseal |
| `VAULT_WATCH_PERIOD` | `60s` | `duration` | Delay between checking Vault seal status |


## Vault commands

To initialize Vault run:
```bash
vault operator init -address=http://127.0.0.1:8200
```