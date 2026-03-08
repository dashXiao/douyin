# TikTok Microservices Demo

一个基于 Go 的 分布式 TikTok 后端微服务示例项目。

## 特点

- 容器化，易于部署
- 实现视频上传、Feed流、点赞评论、关注关系等社交核心功能

## 技术栈

-  Go
-  Hertz
-  Kitex
-  MySQL
-  Redis
-  Kafka
-  etcd

## 快速开始

1. 启动基础依赖环境：

```bash
make env-up
```

2. 将配置写入 etcd：

```bash
docker exec etcd etcdctl put /config/config.yaml "$(cat config/config.yaml)"
```

3. 构建并启动服务：

```bash
make docker
bash docker-run.sh
```

更完整的启动说明见 [start_steps.md](/Users/xyh/GolandProjects/tiktok/start_steps.md)。
