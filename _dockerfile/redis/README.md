### 快速部署


```bash
docker pull redis:3.0-alpine

mkdir -p docker-svc/redis-3.0
cd docker-svc/redis-3.0
touch redis.log && chmod 777 redis.log  # 必需777 否则日志无法写入
mv /path/to/docker-compose.yml .
docker-compose -f docker-compose.yml up -d # `-d` represents start in background

docker ps |grep redis # check status of container
```


### 轻量化部署 redis

先试试水
```shell
docker run --name redis5 --rm -p 6379:6379 redis:5.0 --requirepass "123"
```
如果可以

```shell
docker run --name redis5 -d -p 6379:6379 redis:5.0 --requirepass "123"
```

docker exec -it redis5 redis-cli -a 123