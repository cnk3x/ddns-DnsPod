动态域名解析 DDns for DnsPod

编辑配置文件 ddns.config 格式如下：

{
	"log" : {
		"debug" : false,
		"info" : true,
		"error" : true
	},
	"records" : [{
			"token" : "*****,********************************",
			"host" : "www",
			"domain" : "mydomain.com"
		}, {
			"token" : "*****,********************************",
			"host" : "ddns",
			"domain" : "mydomain.com"
		}
	]
}
