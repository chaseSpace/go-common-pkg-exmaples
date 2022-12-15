FROM docker.io/alpine

# 构建自己的镜像
# docker build -t alpine_test -f alpine_test.dockerfile .
# 测试（经测试，--network=host无法正常使用ssh）
# docker run -it --rm -p 22:22 --hostname=boss alpine_test /bin/bash     # 这种方式进去注意一个问题，只能通过 rc-status -a && rc-service sshd restart 来正常启动ssh，其他方法不行。
# docker run -itd --rm -p 22:22 --hostname=boss alpine_test
# tag
# docker tag alpine_test leigg/alpine_contains_common_cmds:1.0
# docker push leigg/alpine_contains_common_cmds:1.0  // 打开全局代理，避免 net/http: TLS handshake timeout
# docker tag alpine_test leigg/alpine_contains_common_cmds:latest
# docker push leigg/alpine_contains_common_cmds:latest

#RUN echo "#ustc" > /etc/apk/repositories
#RUN echo "https://mirrors.ustc.edu.cn/alpine/v3.6/main/" >> /etc/apk/repositories
#RUN echo "https://mirrors.ustc.edu.cn/alpine/v3.6/community/" >> /etc/apk/repositories

# apk管理
# apk del openssh #删除一个软件
# apk search #查找所以可用软件包
# apk search -v #查找所以可用软件包及其描述内容
# apk search -v 'acf*' #通过软件包名称查找软件包
# apk search -v -d 'docker' #通过描述文件查找特定的软件包

# apk info #列出所有已安装的软件包
# apk info -a zlib #显示完整的软件包信息

# apk upgrade #升级所有软件
# apk upgrade openssh #升级指定软件
# apk upgrade openssh openntp vim   #升级多个软件

# 服务管理（使用rc）
# rc-update add docker boot  开机启动docker
# rc-status -a  列出所有服务
# rc-status -a && rc-service networking restart 重启网络
# rc-status -a && rc-service sshd restart

# 网卡配置文件
# /etc/network/interfaces
# /etc/init.d/networking restart #重启网络
# /etc/init.d/sshd -D

# 解决githubusercontent.com被劫持
# 使用 https://www.ipaddress.com 或者 https://site.ip138.com/ 网站查询 raw.githubusercontent.com 真实IP地址
# 然后将ip填入/etc/hosts文件即可

# 包含系统性能监视命令：uptime top free mpstat iostat
# 内存命令:top free
# 进程命令:ipcs ipcrm lsof strace
# 其他：mount umount df du bash vim curl git

RUN apk update

# 同步上海时间 | 开启ipv4转发
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&  \
    echo "net.ipv4.ip_forward=1" > /etc/sysctl.conf

# 安装基础软件
RUN apk add bash vim curl git openrc strace

# ssh配置 | root用户的密码改为admin
# 进入容器后执行：rc-status -a && rc-service sshd restart  方可启动ssh
# ssh连接错误：WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!
# -- 解决：vi ~/.ssh/known_hosts 删除其中关于目标主机的记录，也可以直接删除这个文件，但会影响你以前连接过的主机
RUN apk add openssh && rc-update add sshd \
        && mkdir /run/openrc/ \
        && touch /run/openrc/softlevel \
        && echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config \
        && echo "PermitRootLogin yes" >> /etc/ssh/sshd_config \
        && ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key \
        && ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key \
        && ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key \
        && ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key \
        && echo "root:admin" | chpasswd

# 开放22端口
EXPOSE 22

# 容器启动时执行ssh启动命令，暂不需要，这样正常启动容器时，无法docker登录容器shell
#CMD ["/usr/sbin/sshd", "-D"]