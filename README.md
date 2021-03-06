# 一键打卡_江苏科技大学版

[![build](https://github.com/yin1999/justhealthreport/actions/workflows/Build.yml/badge.svg)](https://github.com/yin1999/justhealthreport/actions/workflows/Build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/yin1999/justhealthreport)](https://goreportcard.com/report/github.com/yin1999/justhealthreport) [![PkgGoDev](https://pkg.go.dev/badge/github.com/yin1999/justhealthreport.svg)](https://pkg.go.dev/github.com/yin1999/justhealthreport)

## 介绍

项目使用http请求模拟整个打卡过程，速度很快！  
一键打卡，用到就是爽到  
云函数版本请访问[健康打卡_江苏科技大学版_FC](https://github.com/yin1999/justhealthreport_fc)(无服务器版，配置方便，**零成本**)

目前，**最新版本**具有以下特性:

	1. 每日自动打卡
	2. 一次打卡失败，自动重新尝试，可设置最大打卡尝试次数以及重新打卡的等待时间
	3. 日志同步输出到Stderr
	4. 版本查询
	5. 打卡失败邮件通知推送功能(目前支持STARTTLS/TLS端口+PlainAuth登录到SMTP服务器)
	6. 通过环境变量设置http代理(HTTP_PROXY/HTTPS_PROXY均需设置)

### 安装步骤

适用类Unix/windows，想直接使用的，请下载[release](https://gitee.com/allo123/justhealthreport/releases)版本后直接转到[使用说明](#使用说明) 

源码安装依赖[Golang](https://golang.google.cn/)-基于golang开发、[git](https://git-scm.com/)-版本管理工具以及[make](https://www.gnu.org/software/make/)-快速构建，国内使用推荐开启golang的Go module并使用国内的Go proxy服务  
推荐使用[Goproxy.cn](https://goproxy.cn/)或[阿里云 Goproxy](https://developer.aliyun.com/mirror/goproxy)

1. 环境配置，以`Centos`/`Debian`为例

	- 安装Golang[>= 1.16]: [golang.google.cn/doc/install](https://golang.google.cn/doc/install)

	- 安装 git、make:

	   ```bash
	   # Centos
	   sudo yum install git make

	   # Debian/Ubuntu
	   sudo apt install git make
	   ```

2. 通过源码下载、编译

	```bash
	# 配置Goproxy
	go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct  

	# 下载编译
	git clone --depth 1 https://github.com/yin1999/justhealthreport.git
	cd justhealthreport
	make # 若没有安装make，可以使用命令: go run _script/make.go 代替
	```

## 使用说明

### Linux/Windows

1. 授予可执行权限(`源码编译`的可以跳过此步)

	```bash
	chmod +x justhealthreport
	```

2. 运行

	```bash
	# example(set punch time as 9:52)
	./justhealthreport -u username -p password -t 9:52 -save
	# -save: 保存账户信息至文件，第二次启动程序时可不设置用户名、密码两个参数（仅使用: ./healthreport -t 9:52）
	```

	**Linux**用户可使用[systemd](https://systemd.io/)(`recommend`，支持配置开机自启，配置模板：`_script/justhealthreport.service`)或者[screen](https://www.gnu.org/software/screen/)管理打卡进程

### 邮件通知

1. 生成**email.json**

	```bash
	./justhealthreport -g  # 可使用 '-email' 指定配置文件生成目录
	```

2. 修改**email.json**中的的配置，具体说明如下:

	```properties
	to:    收件邮箱(string list)  
	SMTP:  SMTP 配置
	    username:   SMTP用户名(string)
	    password:   SMTP用户密码(string)
	    TLS:        是否为TLS端口(bool)
	    host:       SMTP服务地址(string)
	    port:       SMTP服务端口(需支持STARTTLS/TLS)(int)
	```

3. 重启打卡服务，若提示**Email deliver enabled**，则邮件通知服务已启用
 
### 其它说明

1. 查看版本信息

	```bash
	./justhealthreport -v
	```

2. 帮助信息

	```bash
	./justhealthreport -h
	```

4. 版本更新

	```bash
	git pull
	make
	```
