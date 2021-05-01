# 一键打卡_江苏科技大学版

[![Go Report Card](https://goreportcard.com/badge/gitee.com/allo123/justhealthreport)](https://goreportcard.com/report/gitee.com/allo123/justhealthreport) [![PkgGoDev](https://pkg.go.dev/badge/gitee.com/allo123/justhealthreport)](https://pkg.go.dev/gitee.com/allo123/justhealthreport)

## 介绍
项目使用了headless chrome，感兴趣的可以使用python的Selenium包实现相同的功能。  
一键打卡，用到就是爽到

目前，**最新版本**具有以下特性

    1. 每日自动打卡
    2. 一次打卡失败，自动重新尝试，可设置最大打卡尝试次数以及重新打卡的等待时间
    3. 日志同步输出到Stdout以及log文件

## 安装教程

适用类Unix/windows，推荐在类Unix服务器上运行。

源码安装依赖[golang](https://golang.google.cn/)以及[git](https://git-scm.com/)，国内使用推荐开启golang的Go module并使用国内的Go proxy服务  
推荐使用[goproxy.cn](https://goproxy.cn/)，没有恰饭，没有恰饭，没有恰饭，单纯觉得好用而已。

    git clone https://gitee.com/allo123/one_click2punch.git
    cd one_click2punch_just
    make

下面以Centos为例，实现安装和使用

    #Centos为例
    #环境配置
    yum install -y https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm
    yum install golang screen git -y
    go env -w GO111MODULE=on
    go env -w GOPROXY=https://goproxy.cn,direct

    #下载编译
    git clone https://gitee.com/allo123/justhealthreport.git
    cd justhealthreport
    make

    #运行
    screen ./healthreport
    # 输入必要信息跳出打卡结束信息后
    键盘CTRL+A+D离开screen进程，后台会每天自动运行
    

## 使用说明

使用以下命令获取帮助

    ./healthreport -h


1. linux

先在Linux上安装 Google-Chrome或chromium，百度一下有教程  
然后安装screen

    sudo yum install screen  //centos
    sudo apt-get install screen  //Ubuntu

然后在程序路径下

    #若直接下载了二进制文件，请运行: chmod +x ./healthreport #给予可执行权限
    screen ./healthreport  //服务器请禁用GUI，想看效果可以在PC端上启用GUI

**请使用ctrl+a+d退出screen进程，ctrl+c是用来终止程序的**

2. Windows

依赖chrome，若无法运行程序，提示"没有找到可执行的'chrome'"，请检查chrome是否安装，  
若已经安装了chrome仍然无法正常运行，请将chrome添加到path 
 
然后命令行中执行

    ./healthreport
