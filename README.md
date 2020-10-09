
## Redis-cluster集群创建

<div align=center>
<img src="https://img.shields.io/badge/go-1.15-lightBlue"/>
<img src="https://img.shields.io/badge/redis-5.0.9-brightgreen"/>
<img src="https://img.shields.io/badge/docker-1.18-red"/>
</div>

> 该工具目前适用于装有`docker`的`Linux`系统下



#### 1、使用说明

该工具旨在项目的**开发和前期测试**阶段

虽是为创建集群环境，但该工具只作用于单台机器，通过容器将集群中每个`redis`实例进行隔离，而达到集群效果



#### 2、集群创建

该文件为可执行文件，添加可执行权限

```shell
chmod +x redis-cluster
```

通过`./redis-cluster -h` 来查看相关命令


目前只做了创建`create`相关，删除相关功能还并未做

通过`./redis-cluster create -h` 命令可查看创建需要的选项，分别是 `host`,`ports`,`master`,`replicas`;



`host` 则将作为外界可访问的本机IP地址，如果在云服务器，填写公网IP即可    （必填）

`ports` 则是向外界提供访问`Redis`的端口，用逗号分隔                                        （可不填）

`master` 为集群的master节点个数，理论应为奇数                                              （但是工具将其定为3）

`replicas` 每个master节点个副本数，理论上为任意个                                      （工具将其定位1）



接下来是创建

```
./redis-cluster create -H ${ip}
```

等待其执行完，即可生成6个容器，3主3从，同时集群创建成功。



#### 3、优势

- [x] 操作简便：一行命令即可创建生成集群
- [x] 速度快：整个创建时间在30秒左右
- [x] 硬件要求低：单台机器即可，服务器或者虚拟机都行
- [x] 依赖少：不需要服务器额外安装命令



#### 4、不足

- [ ] 暂时不支持对`master`和`replicas`的自定义赋值
- [ ] Redis的连接没有加入密码验证
- [ ] 暂不支持将集群中数据映射到宿主机
- [ ] 暂不支持window环境
