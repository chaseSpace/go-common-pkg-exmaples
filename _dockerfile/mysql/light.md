## 轻量化部署 mysql (docker on linux)

启动mysql 5.7:

```shell
docker run -d --name mysql \
      -p 3306:3306 \
      -v ~/docker/mysql/data:/var/lib/mysql \
      -v /etc/localtime:/etc/localtime \
      -e MYSQL_ROOT_PASSWORD='123'\
       mysql:5.7
```

启动8.0：

```shell
docker run -d --name mysqlv8 \
      -p 3306:3306 \
      -v ~/docker/mysqlv8/data:/var/lib/mysql \
      -v /etc/localtime:/etc/localtime \
      -e MYSQL_ROOT_PASSWORD='123'\
       mysql:8.0
```

m1 mac需要指定平台拉取镜像：

```shell
docker run -d --name mysql \
      -p 3306:3306 \
      -v ~/docker/mysql/data:/var/lib/mysql \
      -v /etc/localtime:/etc/localtime \
      -e MYSQL_ROOT_PASSWORD='123'\
      --platform linux/x86_64 \
       mysql:5.7
```

启动8.0：

```shell
docker run -d --name mysqlv8 \
      -p 3307:3306 \
      -v ~/docker/mysqlv8/data:/var/lib/mysql \
      -v /etc/localtime:/etc/localtime \
      -e MYSQL_ROOT_PASSWORD='123'\
      --platform linux/x86_64 \
       mysql:8.0
```

其他常用命令：
```shell
# install mysql-client
# https://www.itqaq.com/index/634.html

# 在宿主机尝试连接
mysql -h 127.0.0.1 -u root -p123  # 有时候通过 `-h localhost` 进不去

# 或直接进入容器
docker exec -it mysql mysql -p123
docker exec -it mysqlv8 mysql -p123

# 删除容器
docker stop mysql && docker rm mysql
docker stop mysqlv8 && docker rm mysqlv8
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
CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，自增主键',
    username VARCHAR(50) NOT NULL COMMENT '用户名，最大长度为50个字符，不为空',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '电子邮件地址，最大长度为100个字符，不为空，唯一索引',
    password VARCHAR(255) NOT NULL COMMENT '密码，最大长度为255个字符，不为空',
    gender BIT(1) NOT NULL COMMENT '性别，使用 BIT(1) 类型表示，0 表示女性，1 表示男性，不为空',
    money DECIMAL(2) NOT NULL COMMENT 'money，浮点数',
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期，默认为当前时间戳'
);
# drop table test.users;
```
