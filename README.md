cloudctl
---

### 该项目不够优雅，将进行重构: [https://github.com/skrbox/skrctl](https://github.com/skrbox/skrctl)
### 该项目不够优雅，将进行重构: [https://github.com/skrbox/skrctl](https://github.com/skrbox/skrctl)
### 该项目不够优雅，将进行重构: [https://github.com/skrbox/skrctl](https://github.com/skrbox/skrctl)

> 对于直接部署在操作系统中的服务，我们之前一直沿用supervisord来管理，
> supervisord优点是使用简单、功能强大、便于批量操作，缺点是需要在打包安装时解决依赖问题。
> 因此在dm中优先实践systemd管理服务，之后将作为os中部署进程的首选方式。

> 但是systemd也存在一些问题，如输出不够友好、批量操作不方便等。
> 本程序主要就是对systemctl命令做了一个封装，用来管理通过systemd部署在os中的进程。

## 命令用法

### 命令结构与作用

```text
.
├── cloudctl        # cloudctl可执行文件
├── cloudctl.yml        # 命令配置，指定需要管理的服务
└── services        # 存放若干.service文件，用于批量部署service
```

### 命令用法

```text
[root@192-168-1-99 cloud]# ./cloudctl 
 Usage: ./cloudctl [options] [service_name]
--------------------------------------------------------------------------------
        help [cmd_name]                 to show usage the cmd named                                                                                            
        version                         to show version of the tool
        install                         to install all of service in conf
        install [service_name]          to install the service named
        uninstall                       to uninstall all of service in conf
        uninstall [service_name]        to uninstall the service named
        start                           to start all of service in conf
        start [service_name]            to start the service named
        stop                            to stop all of service in conf
        stop [service_name]             to stop the service named
        restart                         to restart all of service in conf
        restart [service_name]          to restart the service named
        status                          to show status all of service in conf
        status [service_name]           to show status the service named
        enable                          to enable all of service in conf
        enable [service_name]           to enable the service named
        disable                         to disable all of service in conf
        disable [service_name]          to disable the service named
```

- 所有命令都支持单个和批量操作
- install、uninstall用于从systemd安装或卸载service

### 查看单个服务状态

```text
[root@192-168-1-99 ~]# ./cloudctl status phoenix
┌─────────────────┬─────┬─────────┬──────┬───────────┐
│ ServiceName     │ PID │ Status  │ Time │ AutoStart │
├─────────────────┼─────┼─────────┼──────┼───────────┤
│ phoenix.service │ -   │ unknown │ -    │ Yes       │
└─────────────────┴─────┴─────────┴──────┴───────────┘
```

- ServiceName: 服务名称
- PID: 进程主pid
- Status: 进程运行状态
- Time: 进程运行时长
- AutoStart: 是否开机自启
- Version: 服务版本(暂未实现)

### 查看所有服务状态

```text
[root@192-168-1-99 ~]# ./cloudctl status
┌───────────────────────┬──────┬──────────────────┬──────────────┬───────────┐
│ ServiceName           │ PID  │ Status           │ Time         │ AutoStart │
├───────────────────────┼──────┼──────────────────┼──────────────┼───────────┤
│ phoenix.service       │ -    │ unknown          │ -            │ Yes       │
│ prometheus.service    │ 6400 │ active (running) │ 1min 42s ago │ Yes       │
│ dmdb_exporter.service │ -    │ unknown          │ -            │ No        │
│ node_exporter.service │ -    │ unknown          │ -            │ No        │
│ alertmanager.service  │ 6411 │ active (running) │ 1min 42s ago │ Yes       │
└───────────────────────┴──────┴──────────────────┴──────────────┴───────────┘
```

## TODO

- 统一服务version获取方式，status显示版本信息
- 更多功能支持
