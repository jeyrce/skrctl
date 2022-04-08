# skrctl

这是一个封装了systemctl的命令行工具，可用于批量管理systemd部署的服务，直接使用 systemctl 命令的缺点:

- systemctl status 字符繁杂，内容难看
- 无法批量对 service 进行查看和管理

## 使用方法

- 使用 `root` 用户登录 linux 服务器
- 确保 linux 支持 `systemd`
- 需要依赖系统的 `netstat` 命令

### (1) 下载和安装

从源码安装: make 后可得到一个可执行文件 `skrctl`

```shell
git clone https://github.com/skrbox/skrctl.git

cd skrctl

make
```

从github下载二进制: [release](https://github.com/skrbox/skrctl/releases)

### (2) 编写service服务

```text
[Unit]
Description=The nginx HTTP and reverse proxy server
After=network.target remote-fs.target nss-lookup.target

[Service]
Type=forking
PIDFile=/run/nginx.pid
# Nginx will fail to start if /run/nginx.pid already exists but has the wrong
# SELinux context. This might happen when running `nginx -t` from the cmdline.
# https://bugzilla.redhat.com/show_bug.cgi?id=1268621
ExecStartPre=/usr/bin/rm -f /run/nginx.pid
ExecStartPre=/usr/sbin/nginx -t -c /etc/nginx/nginx.conf
ExecStart=/usr/sbin/nginx -c /etc/nginx/nginx.conf
ExecReload=/bin/kill -s HUP $MAINPID
KillSignal=SIGQUIT
TimeoutStopSec=5
KillMode=process
PrivateTmp=true

[Install]
```

如上是一份nginx服务，关于systemd服务写法可以参考：[systemd编写](http://blog.lujianxin.com/x/art/dovhvqvv29g7) ，其中涉及到的服务程序请自行保存在对应位置，然后执行
skrctl 命令开始管理

### (3) 开始管理

```shell
# 假设当前目录为: /tmp
skrctl add /tmp/nginx.service

skrctl start nginx

skrctl ps

# 以下为示例输出
┌───────┬───────┬──────────┬──────────────────┬───────────┬──────┬─────────┐
│ Name  │ PID   │ Port     │ Status           │ Time      │ Auto │ Version │
├───────┼───────┼──────────┼──────────────────┼───────────┼──────┼─────────┤
│ nginx │ 14609 │ 80,60443 │ active (running) │ 26min ago │ Y    │ -       │
└───────┴───────┴──────────┴──────────────────┴───────────┴──────┴─────────┘
```

### (4)直接管理已经启动的service

```shell
[root@10-10-66-26 tmp]# skrctl add /usr/lib/systemd/system/docker.service
[root@10-10-66-26 tmp]# skrctl ps
┌────────┬───────┬──────────┬──────────────────┬────────────┬──────┬─────────┐
│ Name   │ PID   │ Port     │ Status           │ Time       │ Auto │ Version │
├────────┼───────┼──────────┼──────────────────┼────────────┼──────┼─────────┤
│ docker │ 15465 │          │ active (running) │ 2 days ago │ Y    │ -       │
│ nginx  │ 14609 │ 80,60443 │ active (running) │ 30min ago  │ Y    │ -       │
└────────┴───────┴──────────┴──────────────────┴────────────┴──────┴─────────┘
```

## 更多用法

```shell
skrctl --help
skrctl --version
```

## FQA

- 欢迎使用、修改、完善
