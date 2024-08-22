# Commander 命令行工具

## 安装

下载代码库， 进入代码所在目录，执行命令 `go mod tidy`，然后进入 `script` 目录，运行 `build.sh` 脚本

## 使用

`$ go run main.go`
or
`$ ./commander`
```bash
NAME:
   Commander - A new cli application

USAGE:
   Commander [global options] command [command options]

AUTHOR:
   kugouming <kugouming@sina.com>

COMMANDS:
   jsondiff  Compare two http requests json data by field
   listdiff  Compare two http requests json data by list
   stress    Pressure generating tool
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

## 命令说明

### jsondiff

基于多个请求返回的 json 数据做比较，判断字段是否相同。但该工具对返回结果为列表的数据无法进行比较，故请使用 `listdiff` 命令，它会提取列表中数据的关键字段当中唯一 Key，然后比较对应 Key 数据的 md5 值是否相同。

### listdiff

基于多个请求返回的 json 列表数据做比较，它会提取列表中数据的关键字段当中唯一 Key，然后比较对应 Key 数据的 md5 值是否相同。

### stress

这是一个 Go 压力测试工具，支持http和tcp协议。基于 [go-stress-testing](https://github.com/link1st/go-stress-testing) 同步而来，主要是为了方便实用。

## 开发

请操作 cmd 下的 `_template` 目录，创建命令行工具，并将新增工具注册到 `main.go` 文件当中即可。