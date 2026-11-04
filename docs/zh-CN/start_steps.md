# 启动指南

[English](../../start_steps.md) | **中文** · [项目介绍](README.md)

以下命令均在仓库根目录执行。

## 环境要求

- Docker、Docker Compose 插件、Bash 和 Make。
- 可用的主机端口：API 使用 `10001`，MySQL 使用 `3306`，Redis 使用 `6379`，Kafka 使用 `9092`，ZooKeeper 使用 `2181`，etcd 使用 `2379`/`2380`。可选的 Nginx 演示页还需要 `80` 端口。
- 视频功能需要阿里云 OSS 的 endpoint、访问密钥和存储桶。

镜像使用 Go 1.20 构建应用，并包含 FFmpeg。采用 Docker 启动时，无需在主机单独安装 Go。

## 1. 准备配置

复制仓库中的配置模板，以下命令保留了模板文件名的原始拼写：

```bash
cp config/config_exmple.yaml config/config.yaml
mkdir -p config/etcd
```

根据部署环境编辑 `config/config.yaml`。使用仓库脚本启动容器时，各地址配置如下：

| 配置项 | 容器地址或取值 |
| --- | --- |
| `etcd-addr` | `etcd:2379` |
| `etcd.addr` | `etcd:2379` |
| `mysql.addr` | `mysql:3306` |
| `kafka.brokers` | 包含 `kafka:9092` 的 YAML 列表 |
| `redis.addr` | `redis:6379` |
| `redis.password` | 仓库默认 Redis 容器未启用认证，填写 `""`；详见故障排查 |
| `services.api.addr` | 包含 `0.0.0.0:10001` 的 YAML 列表 |
| `services.user.addr` | 包含 `tiktok-user:10002` 的 YAML 列表 |
| `services.chat.addr` | 包含 `tiktok-chat:10003` 的 YAML 列表 |
| `services.follow.addr` | 包含 `tiktok-follow:10004` 的 YAML 列表 |
| `services.interaction.addr` | 包含 `tiktok-interaction:10005` 的 YAML 列表 |
| `services.video.addr` | 包含 `tiktok-video:10006` 的 YAML 列表 |

列表类型的配置需保留模板中的 YAML 列表格式。容器内的 `localhost` 和 `127.0.0.1` 指向当前容器，不能用于连接其他服务容器。

模板中的 MySQL 数据库名、用户名和密码与 Compose 配置一致。请替换 `server.jwt-secret`，并填写 `oss.endpoint`、`oss.accessKey-id`、`oss.accessKey-secret` 和 `oss.bucketname`。视频服务将视频和封面保存到阿里云 OSS，未提供本地存储替代方案。

将完成的配置复制到 etcd 配置监听脚本所挂载的目录：

```bash
cp config/config.yaml config/etcd/config.yaml
```

修改配置时，应同步这两个本地文件。应用容器使用顶层 `etcd-addr` 设置 `ETCD_ADDR`，随后从 etcd 的 `/config/config.yaml` 键读取应用配置。

## 2. 启动基础环境

```bash
bash init.sh
docker compose -p tiktok -f docker-compose.ci.yml up -d
docker compose -p tiktok -f docker-compose.ci.yml ps
```

仓库中的 Compose 文件名为 `docker-compose.ci.yml`，因此需要显式指定。项目名 `tiktok` 会创建 `docker-run.sh` 和 `scripts/run-nginx.sh` 使用的 `tiktok_tiktok` 网络。

等待基础服务初始化完成，推送配置之前先检查 etcd：

```bash
docker exec etcd etcdctl endpoint health
```

## 3. 推送配置到 etcd

配置监听脚本会加载 `config/etcd/config.yaml`。也可以显式推送相同配置，并确认配置键存在：

```bash
docker exec -i etcd etcdctl put /config/config.yaml < config/config.yaml
docker exec etcd etcdctl get /config/config.yaml --keys-only
```

启动应用服务之前，配置必须已经写入 etcd。仅修改本地文件不会配置运行中的服务；应用设置需要存在于 etcd 中。更新配置后，重启受影响的应用容器。

## 4. 构建应用镜像

```bash
make docker
```

此命令构建 `tiktok` 镜像，其中包含 API 网关和五个业务服务。

## 5. 启动应用服务

```bash
bash docker-run.sh
```

脚本会启动 `tiktok-api`、`tiktok-user`、`tiktok-chat`、`tiktok-follow`、`tiktok-interaction` 和 `tiktok-video`，并替换使用该应用镜像的现有容器。

检查 API 网关和视频流：

```bash
curl "http://localhost:10001/ping"
curl "http://localhost:10001/douyin/feed/?latest_time=0"
```

健康检查接口应返回 `{"message":"pong"}`。发布视频之前，Feed 流可能为空。启动异常时，可以查看相关容器日志：

```bash
docker logs tiktok-api
docker logs tiktok-video
```

## 6. 启动 Nginx 和 API 演示页（可选）

```bash
bash scripts/run-nginx.sh
```

在浏览器访问 `http://localhost`。将演示页的 **Base URL** 设为 `http://localhost/douyin`，通过同源 Nginx 代理发送请求。[OpenAPI 文档](tiktok.openapi.json)保留了原有的远程开发服务器地址，因此本地部署时需要修改页面中的地址。

## 故障排查

- **etcd 配置缺失：** 检查 etcd 健康状态及 `/config/config.yaml` 键是否存在，并确认两个本地配置文件内容一致。
- **服务连接失败：** 确认容器均加入 `tiktok_tiktok` 网络、配置使用容器地址，且基础服务已完成启动。
- **Redis 认证：** 默认配置使用 `redis` 镜像，但没有设置 `requirepass`。其中的 `REDIS_PASSWORD` 环境变量不会启用服务端认证。本地默认部署应使用 `redis.password: ""`；如需认证，应同时配置 Redis 服务端和应用密码。
- **Redis 挂载路径：** Compose 文件包含 `/usr/local/redis` 下的绝对主机路径。如果路径不可用或未共享给 Docker，请调整本地 Compose 配置。
- **MySQL 兼容性：** 修改数据之前先查看 `docker logs mysql` 和数据库版本。默认 Compose 文件未挂载 `data/mysql`，清空该目录无法重置 MySQL 容器。
- **视频操作失败：** 检查 OSS 配置及存储桶访问权限。Docker 镜像已包含 FFmpeg；在 Docker 外运行视频服务时，需要单独安装。

### 终端代理

如果终端的 HTTP/HTTPS 代理影响本地 API 请求，可在当前会话中允许本机地址直连：

```bash
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY
export no_proxy="localhost,127.0.0.1${no_proxy:+,$no_proxy}"
export NO_PROXY="localhost,127.0.0.1${NO_PROXY:+,$NO_PROXY}"
```

单次请求也可以显式绕过代理：

```bash
curl --noproxy localhost,127.0.0.1 "http://localhost:10001/ping"
```
