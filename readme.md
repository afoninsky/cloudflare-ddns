# Cloudflare DDNS updater

Personal project for internal purposes so not much docs here. Fetches external IP and updates CloudFlare DNS record (kind of DDNS). There is a bunch of similar project, just decided to write own script is faster than search all of them.

```
kubectl create secret generic -n default coredns-ddns --from-file=.envrc
```
