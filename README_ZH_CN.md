# Autoproxy

[English](./README.md)
[中文](./README_ZH_CN.md) 

该项目提供浏览器http proxy代理服务，支持https、http协议代理；可以部署多级代理；支持TLS协议加密；主要使用与内网主机访问外部网站工具；

## 特性如下
- 支持多种转发方式
- 支持统计控制台
- 支持二次转发
- 支持TLS传输加密
- 支持动态路由
- 支持自定义转发域名设置
- 支持多种平台

## 部署
### 准备提前
- 准备具有公共IP的虚拟云主机
- 准备指定外部TCP协议访问端口，例如：8080

### 使用方式：

下载解压相应平台的软件包；其中包括三个文件；一个autoproxy可执行程序，以及两个参考配置文件；server.yaml 和 client.yaml;

程序命令参数为：

```
Usage of autoproxy.exe:
  -config string
        configure file (default "config.yaml")
  -debug
        enable debug
  -help
        usage help
```

后台启动参考：

`nohup autoproxy -config server.yaml &` 或者 `nohup autoproxy -config client.yaml &`

#### 1、本地代理模式；

在本地代理下面，只需要部署一个autoproxy程序，这个程序作为内网主机访问外网的代理服务；配置参考如下：

```
log:
  path: ./
  filesize: 10485760
  filenumber: 60
local:
  listen: 0.0.0.0:8080
  timeout: 30
  auth:
    - username: user1
      password: uS31k5KLh3NyfvHtFk
    - username: user2
      password: c2O9XJGG0bsJLpt6tr
  mode: local
```

- log: 表示日志记录的目录、单个文件大小、以及文件数量上限；主要是用于审计和问题定位；
- local: 表示程序提供的服务配置，包括监听地址和端口，链路超时时间以及认证方式；如果没有配置认证；则不会进行认证；
- auth: 如果不设置用户名密码，则不进行认证；


然后设置浏览器或者环境变量；
```
export http_proxy="http://用户名:密码@一级代理IP:端口"
export https_proxy="http://用户名:密码@一级代理IP:端口"
```

例如：
```
export http_proxy="http://user1:password1@192.168.3.1:8080"
export https_proxy="http://user1:password1@192.168.3.1:8080"
```

#### 2、本地代理+二级代理模式

二级代理就是在本地代理基础之上，将本地代理的部分或者全部流量通过指定二级代理服务进行转发；可用于复杂的网络环境下，部分网站加速；

本地代理配置：参考压缩包的client.yaml配置文件，默认只需要修改指定二级代理IP地址就可以使用了；

```
log:
  path: ./
  filesize: 10485760
  filenumber: 60
local:
  listen: 0.0.0.0:8080
  timeout: 5
  mode: auto
remote:
  - address: {二级代理IP}:8080
    timeout: 30
    auth:
      username: user1
      password: uS31k5KLh3NyfvHtFk
    tls:
      enable: true
```

- local: 其中 `mode` 有三个选项，分别是：`local`、`auto`、`proxy` ，其中`local` 表示所有流量通过本地路由处理，不会经过二级代理；`auto` 表示根据IP可达性，比如有些本地路由访问不了或者链路超时，则会使用二级代理进行转发，`proxy` 表示所有流量全部经过二级代理；
- remote: 需要访问一个或者多个二级代理的地址，超时时间，认证信息；是否进行TLS加密；如果配置多个二级地址，那么会逐个进行链接尝试；

二级代理部署在VPS侧，需要准备具备一个公网IP的虚拟云主机；

```
log:
  path: ./
  filesize: 10485760
  filenumber: 60
local:
  listen: 0.0.0.0:8080
  timeout: 30
  auth:
    - username: user1
      password: uS31k5KLh3NyfvHtFk
    - username: user2
      password: c2O9XJGG0bsJLpt6tr
  mode: local
  tls:
    enable: true
```

改配置表示二级代理服务端口、认证信息，是否进行TLS加密；如果未配置TLS加密传输，那么一级代理的remote的TLS配置也需要去掉；否则就会链接失败；


#### 3、本地windows UI客户端

本项目提供了小白使用的带UI的客户端，方便使用；在Release 版本下载 `autoproxy_win.zip`
然后解压双击运行即可，本程序是免安装的绿色软件，不会破坏系统；

特性如下：

- 提供基本设置选项
- 转发域名规则
- 远程服务配置
- 最小化和隐藏窗口
- 实时统计控制台
- 本地代理地址和端口设置


![](./docs/main.PNG)

#### 添加二级代理服务

![](./docs/remoteproxy.PNG)

配置完成后，您可以单击“测试”以尝试网络连接性；

#### 编辑域名转发规则

支持几种匹配规则，例如：

- `*.domain.*` : 匹配中间部分域名字段
- `*.domain.com` : 匹配后面域名字段
- `www.domain.*` : 匹配前面域名字段
- `www.domain.com` : 完整匹配域名

![](./docs/domain.PNG)

#### 提供多语言设置

![](./docs/language.PNG)

#### 同步修改本地Internet设置选项

![](./docs/setting.go.PNG)

### 一切就绪，返回主窗口并启动服务； 快乐生活

### [paypal.me](https://paypal.me/lixiangyun)

### Alipay
![](./autoproxy_win/static/sponsor1.jpg)

### Wechat Pay 
![](./autoproxy_win/static/sponsor2.jpg)
