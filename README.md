# 日志文件收集工具

## 接口路径：http://[ip]:[port]/upload-logs
## 请求示例：
  curl --location 'http://XXXX:6088/upload-logs' \
        --form 'file=@"/E:/Pictures/animal03.png"' \
        --form 'desc="{\"type\":\"all\",\"save_path\":\"/data/logs/\"}"'


## MQTT请求示例：
### Topic: event/notice/1806580275824562178
### Message:

```json
    {
      "action": "get_logs",
        "type": "fdm",
        "save_path": "/data/logs/日志"
    }
```
- 1806580275824562178 为 domainId

- 参数说明

| 参数名称 | 类型   | 是否必填  | 说明                                                         |
| -------- | ------ |-------| ------------------------------------------------------------ |
| action   | string | 是     | 业务类型，这里是get_logs                                     |
| type     | string | 否     | 日志类型。检测管理：fdm；网页：bed；检测：fd；所有：all（默认） |
| save_path| string | 否     | 日志存储路径，默认/data/logs/日志                            |



