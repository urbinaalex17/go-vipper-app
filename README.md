# go-vipper-app

## Working with Vault 

Vault must be:
- [Initialized and unsealed](https://learn.hashicorp.com/tutorials/vault/kubernetes-raft-deployment-guide?in=vault/kubernetes#initialize-and-unseal-vault)

Configure Kubernetes authentication:
```bash
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

## Kubernetes Apply

```bash
kubectl apply -f k8s/deployment.yml
```
