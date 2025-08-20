# 日志文件收集工具

## 接口路径：http://[ip]:[port]/upload-logs
## 请求示例：
  curl --location 'http://test.enbo12119.cn:6088/upload-logs' \
        --form 'file=@"/E:/Pictures/animal03.png"' \
        --form 'desc="{\"type\":\"all\",\"save_path\":\"/data/logs/\"}"'


## MQTT请求示例：
### Topic: event/notice/1806580275824562178
### Message:

```json
    {
      "action": "get_logs"
    }
```


