# go-vipper-app

The aim of this repository is to showcase how easily is to work with [Hashicorp Vault](https://www.vaultproject.io) and get specific parameters for your Go applications from your configuration files using [Viper](https://github.com/spf13/viper).

The application just prints the content of the API_KEY and API_SECRET define in the auth.yml configuration file
store in the directory /vault/secrets, retrieved from Vault running withing the Kubernetes Cluster.

```yaml
API_KEY: key
API_SECRET: secret
```

## Working with Vault 

Vault must be:
- [Initialized and unsealed](https://learn.hashicorp.com/tutorials/vault/kubernetes-raft-deployment-guide?in=vault/kubernetes#initialize-and-unseal-vault)

Configure Kubernetes authentication:
```bash
export VAULT_TOKEN=<replace>
vault auth enable kubernetes
```

Configure the Kubernetes authentication method:
```bash
vault write auth/kubernetes/config \
token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
```

Create a secret:
```bash
export VAULT_TOKEN=<replace>
vault secrets enable -path=internal kv-v2
vault kv put internal/vipper/auth key="theKey" secret="theSecret"
```

Verify that the secret is defined:

```bash
vault kv get internal/vipper/auth
```

The output should resemble the following:
```bash
====== Metadata ======
Key              Value
---              -----
created_time     2020-11-07T01:11:11.52715108Z
deletion_time    n/a
destroyed        false
version          1

===== Data =====
Key       Value
---       -----
key       theKey
secret    theSecret
```

The Vault policy:

```bash
vault policy write microservice-vipper - <<EOF
path "internal/data/vipper/auth" {
  capabilities = ["read"]
}
EOF
```

Kubernetes authentication role:

```bash
vault write auth/kubernetes/role/microservice-vipper \
bound_service_account_names=microservice-vipper \
bound_service_account_namespaces=default \
policies=microservice-vipper \
ttl=24h
```

## Kubernetes Deployment

Apply the Kubernetes manifest store in [k8s](k8s) directory:

```bash
kubectl apply -f k8s/deployment.yml
```
Verify that the application is getting the secrets from Vault reading the logs of the vipper container:

```bash
kubectl logs \
$(kubectl get pods -l app=vipper -o jsonpath="{.items[0].metadata.name}") \
-c vipper
```

The output should resemble the following:
```bash
Reading variables using the model..
API_KEY is       theKey
API_SECRET is    theSecret

Reading variables without using the model..
API_KEY is       theKey
API_SECRET is    theSecret
```

## Further documentation and tutorials
- [Viper Go Application](https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63)
- [Injecting Secrets into Kubernetes Pods via Vault Helm Sidecar](https://learn.hashicorp.com/tutorials/vault/kubernetes-sidecar?in=vault/kubernetes)
