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
docker pull mysql:8.0

dst=~/docker-svc/mysql80
data_dir=/mysql_data_dir
container=test-mysql-8.0

mkdir -p $dst
touch $dst/slow-query.log

# 一般root目录空间最大
mkdir -p $data_dir && chmod o+w $data_dir

cd $dst
compose_file=~/docker-compose-mysql80.yml
/bin/cp $compose_file compose.yml

# -d represents start in background
docker-compose -f compose.yml up -d

docker ps |grep mysql # check status of container

# enter mysql shell
docker exec -it $container bash
bash# mysql -uUSER -pPASSWORD

# change root password
mysql> set password for root@localhost = password('new pass'); 
```

```bash
# 开机启动docker-mysql
echo '#!/bin/bash
# chkconfig: 3 88 88
service docker start
sleep 1
docker start test-mysql-8.0
text=$(docker ps --format "table {{.Names}}\t{{.Status}}"|grep mysql)
echo -e $(date +%F_%T) "\t$text" >> /tmp/auto_start.log
' > /etc/init.d/start_docker_mysql

chmod +x /etc/init.d/start_docker_mysql
chkconfig --add start_docker_mysql

```