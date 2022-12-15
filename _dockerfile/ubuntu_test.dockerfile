FROM ubuntu:18.04

# 构建自己的镜像
# docker build -t ubuntu_test -f ubuntu_test.dockerfile .
# 测试
# docker run -it --rm --network=host --hostname=boss_ubuntu ubuntu_test /bin/bash
# docker run -itd --rm --network=host --hostname=boss_ubuntu ubuntu_test
# tag
# docker tag ubuntu_test leigg/ubuntu_contains_common_cmds:1.0
# docker push leigg/ubuntu_contains_common_cmds:1.0   // 打开全局代理，避免 net/http: TLS handshake timeout
# docker tag ubuntu_test leigg/ubuntu_contains_common_cmds:latest
# docker push leigg/ubuntu_contains_common_cmds:latest


# 解决githubusercontent.com被劫持
# 使用 https://www.ipaddress.com 或者 https://site.ip138.com/ 网站查询 raw.githubusercontent.com 真实IP地址
# 然后将ip填入/etc/hosts文件即可

# 包含系统性能监视命令：uptime top free mpstat iostat
# 内存命令:top free
# 进程命令:ipcs ipcrm lsof
# 其他：mount umount df du bash vim curl git

#RUN echo -e "deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse\n\
#deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse\n\
#deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse\n\
#deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse\n\
#deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse\n\
#deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse\n\
#deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse\n\
#deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse\n\
#deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse\n\
#deb-src http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse" > /etc/apt/sources.list

# 这条单独执行，方便构建缓存(修改此文件中已有的命令或调换命令的顺序，就会导致缓存失效)
RUN apt-get update

# 构建时使用本地代理，才能加速软件安装
RUN export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7890

# 同步时区
RUN apt-get install tzdata -y && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 安装基础软件
# net-tools=ifconfig, iputils-ping=ping
RUN apt-get install vim net-tools iputils-ping curl git -y

# 因为网络慢容易下载失败，所以添加新的命令，5个一组的形式下载，这样不影响旧的构建缓存
RUN apt-get install lsof -y