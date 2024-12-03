## 介绍

网上对JWT伪造的工具较多，但是go语言常见的session库储存在cookie时没有采取JWT形式，造成了一定麻烦。笔者分享一个开发的gorilla/session伪造工具，适用于如下场景。

```go
var store = sessions.NewCookieStore([]byte("xxxxxxx"))

func Index(c *gin.Context) {
	session, err := store.Get(c.Request, "session-xxx")
```

## 使用方式

首先在keys.txt填入你猜测的密钥，一行一个。

然后运行如下命令。

```sh
go run bruteForce.go
```

导出的outs.txt文件可以在burpsuit的爆破模块里使用。

## 温馨提示

由于go的特性，对空字符串这类弱密钥不友好，其余弱密码没有问题。

## 贡献

欢迎任何人继续升级本项目。