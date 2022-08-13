## klcc-tools

### 说明

```
调用minio的SDK实现在 data(config.yaml 可更改需要上传的路径) 下的目录为bucket名称，并设置为 * RW
```

### 使用

```
git clone https://github.com/klcc-c/klcc-tools.git
cd klcc-tools
go mod tidy
go build main.go
./klcc-tools upfile
```

TODO

```
创建配置文件中描述的数据库，并执行sql文件恢复数据库
```

