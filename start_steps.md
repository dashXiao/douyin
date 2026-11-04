# Startup Guide

**English** | [中文](docs/zh-CN/start_steps.md) · [Project overview](README.md)

Run the commands below from the repository root.

## Prerequisites

- Docker with the Docker Compose plugin, Bash, and Make.
- Available host ports: `10001` for the API, `3306` for MySQL, `6379` for Redis, `9092` for Kafka, `2181` for ZooKeeper, and `2379`/`2380` for etcd. The optional Nginx demo also uses port `80`.
- An Alibaba Cloud OSS endpoint, access keys, and bucket for video operations.

The image builds the application with Go 1.20 and includes FFmpeg. A separate host Go installation is not required for the Docker workflow.

## 1. Prepare the Configuration

Copy the checked-in template; the filename below matches the repository's spelling:

```bash
cp config/config_exmple.yaml config/config.yaml
mkdir -p config/etcd
```

Edit `config/config.yaml` for your environment. For containers launched by the included scripts, use these addresses:

| Setting | Container address or value |
| --- | --- |
| `etcd-addr` | `etcd:2379` |
| `etcd.addr` | `etcd:2379` |
| `mysql.addr` | `mysql:3306` |
| `kafka.brokers` | A YAML list containing `kafka:9092` |
| `redis.addr` | `redis:6379` |
| `redis.password` | `""` for the supplied unauthenticated Redis container; see troubleshooting |
| `services.api.addr` | A YAML list containing `0.0.0.0:10001` |
| `services.user.addr` | A YAML list containing `tiktok-user:10002` |
| `services.chat.addr` | A YAML list containing `tiktok-chat:10003` |
| `services.follow.addr` | A YAML list containing `tiktok-follow:10004` |
| `services.interaction.addr` | A YAML list containing `tiktok-interaction:10005` |
| `services.video.addr` | A YAML list containing `tiktok-video:10006` |

Keep list settings in the template's YAML list format. `localhost` and `127.0.0.1` refer to the current container, so they cannot address a different service container.

The template's MySQL database, username, and password match the supplied Compose configuration. Replace `server.jwt-secret` with your own value, and fill in `oss.endpoint`, `oss.accessKey-id`, `oss.accessKey-secret`, and `oss.bucketname`. The video service uses Alibaba Cloud OSS for videos and covers; it does not provide a local-storage fallback.

Copy the completed configuration to the directory mounted by the etcd configuration watcher:

```bash
cp config/config.yaml config/etcd/config.yaml
```

Keep both local copies synchronized when changing configuration. Application containers use the top-level `etcd-addr` to set `ETCD_ADDR`, then load the application settings from the etcd key `/config/config.yaml`.

## 2. Start the Infrastructure

```bash
bash init.sh
docker compose -p tiktok -f docker-compose.ci.yml up -d
docker compose -p tiktok -f docker-compose.ci.yml ps
```

The repository's Compose file is named `docker-compose.ci.yml`, so specify it explicitly. The project name `tiktok` creates the `tiktok_tiktok` network expected by `docker-run.sh` and `scripts/run-nginx.sh`.

Wait for the dependencies to finish initializing. Check etcd before publishing the configuration:

```bash
docker exec etcd etcdctl endpoint health
```

## 3. Publish the Configuration to etcd

The watcher loads `config/etcd/config.yaml`. You can explicitly publish the same configuration and confirm that the key exists:

```bash
docker exec -i etcd etcdctl put /config/config.yaml < config/config.yaml
docker exec etcd etcdctl get /config/config.yaml --keys-only
```

Publish the configuration before starting the application services. Editing a local file alone does not configure a running service; the settings must be available in etcd. Restart the affected application containers after updating the configuration.

## 4. Build the Application Image

```bash
make docker
```

This builds the `tiktok` image containing the API gateway and all five business services.

## 5. Start the Application Services

```bash
bash docker-run.sh
```

The script starts `tiktok-api`, `tiktok-user`, `tiktok-chat`, `tiktok-follow`, `tiktok-interaction`, and `tiktok-video`. It replaces existing containers that use the application image.

Check the API gateway and video feed:

```bash
curl "http://localhost:10001/ping"
curl "http://localhost:10001/douyin/feed/?latest_time=0"
```

The health endpoint should return `{"message":"pong"}`. The feed may be empty until videos are published. For startup failures, inspect the relevant container:

```bash
docker logs tiktok-api
docker logs tiktok-video
```

## 6. Start Nginx and the API Playground (Optional)

```bash
bash scripts/run-nginx.sh
```

Open `http://localhost` in your browser. Set the playground's **Base URL** to `http://localhost/douyin` to send requests through the same-origin Nginx proxy. The [OpenAPI specification](demo/tiktok.openapi.json) retains the original remote development-server URL, so update the field for your local deployment.

## Troubleshooting

- **Missing etcd configuration:** check etcd health and confirm `/config/config.yaml` exists. Ensure the two local configuration files contain the same settings.
- **Services cannot connect:** confirm that all containers share `tiktok_tiktok`, that the configuration uses container addresses, and that dependencies have finished starting.
- **Redis authentication:** the supplied service uses the `redis` image without a `requirepass` setting. Its `REDIS_PASSWORD` environment variable does not configure server authentication. Use `redis.password: ""` for this local setup, or configure Redis authentication and match the application password.
- **Redis mounts:** the Compose file includes absolute host paths under `/usr/local/redis`. Adjust those paths in your local Compose configuration if they are unavailable or not shared with Docker.
- **MySQL compatibility:** inspect `docker logs mysql` and the database version before changing stored data. The supplied Compose file does not mount `data/mysql`, so clearing that directory does not reset the MySQL container.
- **Video operations fail:** verify the OSS settings and bucket access. FFmpeg is included in the Docker image and must be installed separately when running the video service outside Docker.

### Terminal Proxies

If a terminal HTTP/HTTPS proxy interferes with local API requests, allow direct connections to local addresses for the current session:

```bash
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY
export no_proxy="localhost,127.0.0.1${no_proxy:+,$no_proxy}"
export NO_PROXY="localhost,127.0.0.1${NO_PROXY:+,$NO_PROXY}"
```

For a single request, bypass the proxy explicitly:

```bash
curl --noproxy localhost,127.0.0.1 "http://localhost:10001/ping"
```
