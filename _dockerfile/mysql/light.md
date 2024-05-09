## 轻量化部署 mysql (docker on linux)

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
```

其他常用命令：
```shell
# 在宿主机尝试连接
mysql -h 127.0.0.1 -u root -p123  # 有时候通过 `-h localhost` 进不去

# 或直接进入容器
docker exec -it mysql mysql -p123

# 删除容器
docker stop mysql && docker rm mysql
```

## 设置mysql远程登录
默认mysql仅支持本机访问，所以如果需要，则通过上面的常用命令进入mysql shell，执行下面指令：
```shell
# 123是root密码
GRANT ALL PRIVILEGES ON *.* TO root@'%' IDENTIFIED BY '123' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

## 测试表

```
create database test;
CREATE TABLE test.users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    gender BIT(1) NOT NULL,
    money DECIMAL(2) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
# drop table test.users;
```
