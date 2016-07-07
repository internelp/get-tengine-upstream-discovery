# get-tengine-upstream-discovery
基于get-tengine-upstream，自动发现的数据采集工具。


##发现规则
nginx.upstream.discovery

####相关值
- #SRVNAME		后端服务器名称（ip和端口）
- #STATUS		后端服务器状态（up=1，down=0）
	返回数据如下，负载均衡池中有几个服务器，就会出现几个。
```json
{
  "data": [
    {
      "{#SRVNAME}": "127.0.0.1:7070"
    },
    {
      "{#SRVNAME}": "127.0.0.1:9090"
    }
  ]
}
```
