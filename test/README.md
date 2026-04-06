1. Create a new secret:
```
k create secret generic nats-secret --from-literal=token=$(k get secret -o yaml eventbus-default-client | yq .data.client-auth | base64 -d | yq .token)
```
2. Apply manifest
```
k apply -f eventsource.yaml
```
