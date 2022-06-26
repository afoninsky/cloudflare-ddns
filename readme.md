# Cloudflare DDNS updater

Personal project for internal purposes so not much docs here. Fetches external IP and updates CloudFlare DNS record (kind of DDNS). There is a bunch of similar project, just decided to write own script is faster than search all of them.

```
echo CLOUDFLARE_API_TOKEN=$CLOUDFLARE_API_TOKEN > /tmp/env
echo CLOUDFLARE_ZONE=$CLOUDFLARE_ZONE >> /tmp/env
kubectl create secret generic -n default cloudflare --from-env-file=/tmp/env
kubectl create job --from=cronjob/ddns ddns-test
```
