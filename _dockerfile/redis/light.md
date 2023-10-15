### 轻量化部署 redis

先试试水
```shell
docker run --name redis5 --rm -p 6379:6379 redis:5.0 --requirepass "123"
```
如果可以，再运行

```shell
docker run --name redis5 -d -p 6379:6379 redis:5.0 --requirepass "123"
```

docker exec -it redis5 redis-cli -a 123