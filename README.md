# whoamimcp


### Environment variables

| Name | Default | Description |
| ---- | ------- | ----------- |
| `WHOAMI_NAME` | `whoamimcp` | Name tag |
| `WHOAMI_HOST` | `0.0.0.0` | Host to listen on |
| `WHOAMI_PORT_NUMBER` | `80` | Port to listen on |


### Arguments

Takes precedence over environment variables

| Name | Default | Description |
| ---- | ------- | ----------- |
| `name` | `whoamimcp` | Name tag |
| `addr` | `0.0.0.0:80` | Address to listen on |


### Features

| Type | Name | Description |
| ---- | ---- | ----------- |
| Prompt | greet | Generates simple greet prompt |
| Tool | greet | Replies with simple greet text |
| Tool | whoami | Replies with server information |
| Tool | raise_error | Simulates error |


### Run

#### Docker

```bash
# default
docker run --rm -it -n whoamimcp traefik/whoamimcp:latest
```

```bash
# configured
docker run --rm -it \
    -n whoamimcp \
    -e "TRAEFIK_NAME=whoamimcp" \
    -e "TRAEIFK_HOST=0.0.0.0" \
    -e "TRAEFIK_PORT_NUMBER=80" \
    traefik/whoamimcp:latest
```