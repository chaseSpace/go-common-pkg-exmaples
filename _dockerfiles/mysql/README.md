### 常见操作
```bash
docker exec -it test-mysql-5.7 bash  // 进入容器shell
docker logs test-mysql-5.7  // 查看日志

# dump db
$ docker exec some-mysql sh -c 'exec mysqldump --all-databases -uroot -p"$MYSQL_ROOT_PASSWORD"' > /some/path/on/your/host/all-databases.sql

# restore db
$ docker exec -i some-mysql sh -c 'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD"' < /some/path/on/your/host/all-databases.sq

```

### 快速部署

```bash
docker pull mysql:5.7

mkdir -p docker-svc/mysql-5.7
cd docker-svc/mysql-5.7
mv /path/to/docker-compose.yml .
docker-compose -f docker-compose.yml up -d # `-d` represents start in background

docker ps |grep mysql # check status of container
```