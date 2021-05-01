# 一键打卡_江苏科技大学版

[![build](https://github.com/yin1999/justhealthreport/actions/workflows/Build.yml/badge.svg)](https://github.com/yin1999/justhealthreport/actions/workflows/Build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/yin1999/justhealthreport)](https://goreportcard.com/report/github.com/yin1999/justhealthreport) [![PkgGoDev](https://pkg.go.dev/badge/github.com/yin1999/justhealthreport.svg)](https://pkg.go.dev/github.com/yin1999/justhealthreport)

## 介绍

项目使用http请求模拟整个打卡过程，速度很快！  
一键打卡，用到就是爽到

目前，**最新版本**具有以下特性

    1. 每日自动打卡
    2. 一次打卡失败，自动重新尝试，可设置最大打卡尝试次数以及重新打卡的等待时间
    3. 日志同步输出到Stdout以及log文件
    4. 版本查询
    5. 打卡失败邮件通知推送功能(目前支持STARTTLS/TLS端口+PlainAuth登录到SMTP服务器)

### 安装步骤

1. 环境配置，以CentOS 7为例

    安装软件：Golang[>= 1.16]、screen、git、make

       yum install -y golang screen git make

2. 通过源码下载、编译

       go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct  #配置Goproxy

       #下载编译
       git clone https://github.com/yin1999/justhealthreport.git
       cd justhealthreport
       make

3. 运行

       screen ./justhealthreport
       # 输入必要信息跳出验证账号密码成功后
       # 键盘CTRL+A+D离开screen进程，后台会每天自动运行

## 使用说明

### linux

1. 安装screen

       sudo yum install screen  #CentOS
       sudo apt-get install screen  #Ubuntu

2. 运行

       chmod +x justhealthreport # 直接下载发行版时需要授予可执行权限
       screen ./justhealthreport

**请使用ctrl+a+d退出screen进程，ctrl+c是用来终止程序的**

### Windows

命令行中执行

    .\justhealthreport

### 邮件通知

1. 复制**email-template.json**命名为**email.json**

       cp email-template.json email.json  #Linux命令

2. 修改**email.json**的权限为仅使用者可读、可修改(0600权限)，保证数据安全，Windows用户请通过文件属性删除**Users**的读取权限

       chmod 600 email.json  #Linux命令

3. 修改**email.json**中的的配置，具体说明如下:

       to:    收件邮箱(string list)  
       SMTP:  SMTP 配置  
         username:   SMTP用户名(string)  
         password:   SMTP用户密码(string)  
         TLS:        是否为TLS端口(bool)  
         host:       SMTP服务地址(string)  
         port:       SMTP服务端口(需支持STARTTLS/TLS)(int)

4. 重启打卡服务，若提示**Email deliver enabled**，则邮件通知服务已启用
 
### 其它说明

1. 查看版本信息

       ./justhealthreport -v

2. 帮助信息

       ./justhealthreport -h

3. 验证SMTP服务

       ./justhealthreport -c

4. 版本更新

       git pull
       make clean && make
