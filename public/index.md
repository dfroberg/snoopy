# Welcome to Snoopy
This is the Helm Chart Repo for Snoopy.

Snoopy subscribes to events on the Ethereum network you specify and
spits out stats and information about blocks it has gathered since
it started.
---
How to get it;
~~~
helm repo add snoopy https://dfroberg.github.io/snoopy/
helm repo update
~~~
Take a look;
~~~
helm search repo snoopy
~~~
How to install it;
~~~
helm upgrade snoopy snoopy/snoopy \
      --install \
      --namespace snoopy \
      --create-namespace \
      --wait \
      --set snoopy.image.tag=latest \
      --set snoopy.ingress.enabled=true \
      --set snoopy.ingress.domain.prefix="" \
      --set snoopy.ingress.domain.base=snoopy.local \
      --set snoopy.metrics.enabled=true
~~~
How to test it;
~~~
helm test snoopy --namespace snoopy
~~~

**Homepage:** <https://github.com/dfroberg/snoopy>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| dfroberg | <danny@froberg.org> |  |

## Source Code

* <https://github.com/dfroberg/snoopy/tree/master/snoopy>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| common | object | `{"snoopyApiToken":"TestToken"}` | Common values for all services |
| common.snoopyApiToken | string | `"TestToken"` | This is optional, will be pupulated by a random string if not defined or already present in a secret. |
| snoopy | object | `{"env":[{"name":"TZ","value":"Europe/Stockholm"}],"image":{"pullPolicy":"Always","repository":"dfroberg/pleo-snoopy","tag":"latest"},"ingress":{"annotations":{},"domain":{"base":"snoopy.local","prefix":"","suffix":""},"enabled":true,"ingressClassName":"traefik","labels":{}},"resources":{"limits":{"memory":"4096Mi"},"requests":{"memory":"4096Mi"}},"service":{"port":8000},"metrics":{"enabled":true,"port":8000}}` | Values for snoopy service |
| snoopy.env | list | `[{"name":"TZ","value":"Europe/Stockholm"}]` | Environment vars to set |
| snoopy.ingress.enabled | bool | `true` | Enable ingress |
| snoopy.ingress.annotations | object | `{}` | Ingress annotations |
| snoopy.ingress.labels | object | `{}` | Ingress labels |
| snoopy.ingress.ingressClassName | string | `"traefik"` | IngressClassname |
| snoopy.ingress.domain | object | `{"base":"snoopy.local","prefix":"","suffix":""}` | Build host string |
| snoopy.service.port | int | `9080` | Port number (Defaults to 8000) |
| snoopy.metrics.enabled | bool | `true` | Enable if you wish to enable prometheus metrics |
| snoopy.metrics.port | int | `2112` | Port number (Defaults to 8000) |
| snoopy.resources | object | `{"limits":{"memory":"1024Mi"},"requests":{"memory":"1024Mi"}}` | Resource limits |


