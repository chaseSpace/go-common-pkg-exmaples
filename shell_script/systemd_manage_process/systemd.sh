<<EOF
使用systemd来管理服务进程

使用：
./systemd.sh your_service_file your_bin_file

This will do:
	/bin/cp your_bin_file $BIN_STORE_DIR -f
	/bin/cp your_service_file $SYSTEMD_SVC_FILE_DIR -f

执行后如何check/start/stop
sl status svc.target
sl stop svc.target
sl start svc.target
sl unmask svc.target   禁用服务，禁用后enable/start都不能生效

建议：systemd仍然不适合管理单个服务的多进程部署，在扩/减进程数量时需要自己使用脚本实现；
本脚本的初衷是尝试探究systemd在管理单个服务的多进程部署能力，目前看来目标可以勉强实现，但管理仍不太方便，所以还是建议使用supervisor

<<systemd擅长的领域是单进程管理>>
EOF
alias sl=systemctl

service_f=$1
bin_f=$2

BIN_STORE_DIR='/usr/local/sbin/'
SYSTEMD_SVC_FILE_DIR='/run/systemd/system/'

# need clean some files inside these path, before start svc
# 系统自动复制link到此路径，非人为控制，我们需要在每次部署时删除这个目录下对应的target文件（在缩容时清除多余文件）
# 但是！！！不能这样去做，因为这个路径下还有系统的服务文件，很有可能出现误删除，由于我们也无法修改这个路径
# 所以只能放弃这种做法！
#AUTO_COPY_SYSLINK_PATH='/etc/systemd/system/'

if [[ ! -f $service_f ]]; then
	echo 'you must provide a valid service file path'
	exit 1
elif [[ ! -x $bin_f ]]; then
	echo "Not a executable file, execute> chmod +x $bin_f"
	chmod +x $bin_f

	if [[ ! $? == 0 ]]; then
		exit 1
	fi
fi

# do clean


# svc@.service ==> svc@
f_name=${service_f%.*}
pure_f_name=${f_name::-1}

# copy binary file
/bin/cp $bin_f $BIN_STORE_DIR -f

target=$pure_f_name.target

echo '[Install]
WantedBy=multi-user.target' >$SYSTEMD_SVC_FILE_DIR/$target

# copy .service file
/bin/cp $service_f $SYSTEMD_SVC_FILE_DIR

systemctl daemon-reload
systemctl enable $pure_f_name@{1..5}.service  # svc@{1..N}.service
# Created symlink from /etc/systemd/system/simple_http2.target.wants/simple_http@1.service to
# /run/systemd/system/simple_http@.service.
systemctl restart $target
#systemctl enable target  # start on boot

echo ">systemctl | grep $pure_f_name@"
systemctl | grep ${pure_f_name}'@'

echo ''

echo ">netstat -antlp|grep $pure_f_name"
netstat -antlp|grep $pure_f_name
