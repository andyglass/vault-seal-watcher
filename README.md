# Vault Seal Watcher

> Dockerhub repository https://hub.docker.com/r/krezz/vault-seal-watcher

* [Usage](#usage)
  - [Docker](#docker)
  - [Kubernetes](#kubernetes)
* [Environment variables](#environment-variables)


<a name="usage"></a>
## Usage

<a name="docker"></a>
### Docker

To start docker image run:
```bash
docker run --rm \
    -e VAULT_ADDR="http://localhost:8200"
    -e VAULT_UNSEAL_KEYS="5g109m1QWGTNCBEpSGVsx4QwxjdxrIKYrUVUkOndTWSp,Dq3kp4KQXw/pPizGwWMF9OC2pFEAwBI5IPiWd/WlQ5V/"
    -p 8000:8000
    krezz/vault-seal-watcher
```

Check vault-seal-watcher liveness:
```bash
curl -Ss http://localhost:8000/health | jq
```
```json
{
  "health": "ok"
}
```

<a name="kubernetes"></a>
### Kubernetes with Hashicorp Vault Helm chart

Initialize Vault server instance:
```bash
vault operator init -address=http://127.0.0.1:8200
```

Create kubernetes secret with comma-separated list of Vault unseal keys (base64 encoded):
```yaml
apiVersion: v1
kind: Secret
metadata:
name: vault-unseal-keys
namespace: vault
labels:
    app.kubernetes.io/name: vault-unseal-keys
    app.kubernetes.io/part-of: vault
type: kubernetes.io/Opaque
data:
  VAULT_UNSEAL_KEYS: |
    VjFFOCs1dk9HR01HVmdscGlYTVdiMXFIdlV6cXJLSTJsWjJtT0N3QlNoV2IsRVdZT0pVd0Q3aUQ1
    SERFcU5aTVRlTWgvWUZUMmY2QjJMOVhJNmFQakg4NjYsTGdDeU80djJCV0ZDYzhCem5oZmd1K1JY
    YjVaWTFRZXJNQU9pNy9NdkFYamQsL0lXTkRsNjRmSEw3RFZScVRob0NycGdjTkhIaUc0NzEzM0ZF
    WlBrc2NsK04sbkxESHVuTWsrNStJZFAxak9NSE1uSG1INWVtWVVuWStWUzNMR0d6bUVJWTgK
```

Add Vault-seal-watcher as extra container to official Vault Helm chart values file:
```yaml
server:
  extraContainers:
  - name: vault-seal-watcher
    image: krezz/vault-seal-watcher:latest
    env:
    - name: VAULT_UNSEAL_KEYS
      valueFrom:
        secretKeyRef:
          name: vault-unseal-keys
          key: VAULT_UNSEAL_KEYS
    ports:
    - name: http
      protocol: TCP
      containerPort: 8000
    resources:
      limits:
        cpu: 100m
        memory: 128Mi
      requests:
        cpu: 50m
        memory: 64Mi
    livenessProbe:
      tcpSocket:
        port: http
      initialDelaySeconds: 5
      periodSeconds: 5
      timeoutSeconds: 5
      successThreshold: 1
      failureThreshold: 3
    readinessProbe:
      httpGet:
        scheme: HTTP
        port: http
        path: /health
      initialDelaySeconds: 5
      periodSeconds: 5
      timeoutSeconds: 5
      successThreshold: 1
      failureThreshold: 1
```


<a name="environment-variables"></a>
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
| `VAULT_SKIP_TLS_VERIFY` | `false` | `bool` | Skip TLS verify on requests to Vault server |
| `VAULT_UNSEAL_KEYS` | `nil` | `string` | Comma-separated list of key-shares |
| `VAULT_UNSEAL_DELAY` | `200ms` | `duration` | Delay between using key-shares to unseal |
| `VAULT_WATCH_PERIOD` | `60s` | `duration` | Delay between checking Vault seal status |
