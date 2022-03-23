### 常见操作
```bash
docker exec -it test-mysql-8.0.28 bash  // 进入容器shell
docker logs test-mysql-8.0.28  // 查看日志

# dump db
$ docker exec some-mysql sh -c 'exec mysqldump --all-databases -uroot -p"$MYSQL_ROOT_PASSWORD"' > /some/path/on/your/host/all-databases.sql

# restore db
$ docker exec -i some-mysql sh -c 'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD"' < /some/path/on/your/host/all-databases.sq

```

### 快速部署

```bash
version=8.0.28
version2=${version::3} # 5.7
version3=${version:0:1}${version:2:1} # 57
echo mysql-- $version2 $version3

docker pull mysql:$version

mysql_dir=~/docker_mysql80/
data_dir=/$mysql_dir/data

# define container var
c=test-mysql-$version

mkdir -p $mysql_dir 
mkdir -p $mysql_dir/conf.d  # 把config文件放进去, 文件名必须是my.cnf
mkdir -p $mysql_dir/log/ && chmod 777 $mysql_dir/log/
mkdir -p $data_dir && && chmod 777 $data_dir

cd $mysql_dir
# 一般root目录空间最大

# 安装docker
# 。。。
# 安装docker-compose
curl -L https://get.daocloud.io/docker/compose/releases/download/1.26.2/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
docker-compose --version

# 创建docker-compose.yml, 将此目录下的docker-compose-mysql80.yml内容复制进去
vi docker-compose.yml
# 启动（完全启动mysql需要30s-1min）
docker-compose up -d
# 注意：启动后的mysql root只能本地登录
docker ps |grep mysql # check status of c

# enter mysql shell
docker exec -it $c bash
bash# mysql -uroot -phoulangfeilang

# 远程登录配置，change root password(8.0)
#use mysql;
#create user root@'%' identified by 'houlangfeilang';
#grant all privileges on *.* to root@'%' with grant option;
# 执行第一条
update mysql.user set authentication_string='',Host='%' where user='root';
# 执行第二条，停顿1s 再执行下一条，否则可能会执行失败： ERROR 1396 (HY000)
ALTER USER root@'%' IDENTIFIED WITH mysql_native_password BY 'houlangfeilang';

# 开机启动docker-mysql
echo '#!/bin/bash 
# chkconfig: 3 88 88 
service docker start 
sleep 1 
docker start test-mysql-8.0.28 
text=$(docker ps --format "table {{.Names}}\t{{.Status}}"|grep mysql) 
echo -e $(date +%F_%T) "\t$text" >> /tmp/auto_start.log 
' > /etc/init.d/start_docker_mysql80

chmod +x /etc/init.d/start_docker_mysql80
chkconfig --add start_docker_mysql80
```