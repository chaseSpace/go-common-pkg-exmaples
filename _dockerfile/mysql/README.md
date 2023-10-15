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
version=5.7.35
version2=${version::3} # 5.7
version3=${version:0:1}${version:2:1} # 57
echo mysql-- $version2 $version3

docker pull mysql:$version

dst=~/docker-svc/mysql$version3
data_dir=/mysql$version3_data_dir
container=test-mysql-$version2

mkdir -p $dst
touch $dst/slow-query.log

# 一般root目录空间最大
mkdir -p $data_dir && chmod o+w $data_dir

cd $dst
compose_file=~/docker-compose-mysql$version3.yml
/bin/cp $compose_file compose.yml

# 安装docker
# 。。。
# 安装docker-compose
curl -L https://get.daocloud.io/docker/compose/releases/download/1.26.2/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
docker-compose --version

# -d represents start in background
docker-compose -f compose.yml up -d
# 注意：启动后的mysql root只能本地登录

docker ps |grep mysql # check status of container

# enter mysql shell
docker exec -it $container bash
bash# mysql -uUSER -pPASSWORD

# change root password
mysql> set password for root@localhost = password('new pass'); 

# 开机启动docker-mysql
echo '#!/bin/bash 
# chkconfig: 3 88 88 
service docker start 
sleep 1 
docker start test-mysql-8.0 
text=$(docker ps --format "table {{.Names}}\t{{.Status}}"|grep mysql) 
echo -e $(date +%F_%T) "\t$text" >> /tmp/auto_start.log 
' > /etc/init.d/start_docker_mysql80

chmod +x /etc/init.d/start_docker_mysql80
chkconfig --add start_docker_mysql80

echo '#!/bin/bash 
# chkconfig: 3 88 88 
service docker start 
sleep 1 
docker start test-mysql-5.7
text=$(docker ps --format "table {{.Names}}\t{{.Status}}"|grep mysql) 
echo -e $(date +%F_%T) "\t$text" >> /tmp/auto_start.log 
' > /etc/init.d/start_docker_mysql57
chmod +x /etc/init.d/start_docker_mysql57
chkconfig --add start_docker_mysql57
```