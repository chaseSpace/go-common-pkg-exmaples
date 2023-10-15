### 轻量化部署 mysql (docker on linux)

先尝试运行看是否报错：
```shell
docker run --rm --name mysql \
      --network host \
      -v ~/docker/mysql/data:/var/lib/mysql \
      -v /etc/localtime:/etc/localtime \
      -e MYSQL_ROOT_PASSWORD='123'\
       mysql  # mac上替换为 mariadb
       
docker exec -it mysql mysql -p123
mysql -h 127.0.0.1 -u root -p123  # mariadb也需要用 mariadb client连接
```

如果没报错，就crtl+C退出，运行下面命令：
```shell
docker run -d --name mysql \
      --network host \
      -v ~/docker/mysql/data:/var/lib/mysql \
      -v /etc/localtime:/etc/localtime \
      -e MYSQL_ROOT_PASSWORD='123'\
       mysql:5.7  # mac上替换为 mariadb:5.7
mysql -h 127.0.0.1 -u root -p123  # 有时候通过 `-h localhost` 进不去
docker stop mysql && docker rm mysql

docker exec -it mysql mysql -p123

# 远程登录（但上述命令指定了容器网络是host，所以宿主机登录可以直接认为是本机登录）
GRANT ALL PRIVILEGES ON *.* TO root@'%' IDENTIFIED BY '123' WITH GRANT OPTION; 
FLUSH PRIVILEGES;
```