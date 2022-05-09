# 代码生成器
依赖于项目 github.com/armory-toolkit/armory-web-starter-go

## 使用方法
+ 在根目录中创建 main 编辑如下
```golang
//  go:generate go env -w GOPRIVATE=github.com/armory-toolkit,bywin.com
//  go:generate go env -w GOINSECURE=github.com/armory-toolkit,bywin.com
//go:generate go install golang.org/x/tools/cmd/goimports@latest
//go:generate go install github.com/justxi1/bywin_generate@master
//go:generate go get -u github.com/armory-toolkit/armory-web-starter-go@dev
//go:generate bywin_generate --type=mysql --override --projectName=generate_localTest --ip=xxx.xxx.xxx.138 --port=32556 --schema=test --user=root --password=root
//go:generate go generate ./cmd/app/main.go
//go:generate go mod tidy
//go:generate goimports -w .
//go:generate gofmt -l -w -s .
package main

func main() {

}
```
+ 打开终端，在 main 目录下执行
```shell
go generate
```

## 说明
```shell
// 配置私库
//  go:generate go env -w GOPRIVATE=github.com/armory-toolkit
//  go:generate go env -w GOINSECURE=github.com/armory-toolkit

// 下载 go官方工具 goimports
//go:generate go install golang.org/x/tools/cmd/goimports@latest

// 下载项目工具 bywin_generate
//go:generate go install github.com/justxi1/bywin_generate@master

// 获取依赖
//go:generate go get -u github.com/armory-toolkit/armory-web-starter-go@dev

// 连接数据库生成代码
//go:generate bywin_generate --type=mysql --override --projectName=generate_localTest --ip=192.168.96.138 --port=32556 --schema=test --user=root --password=root

// 生成 swagger 文档
//go:generate go generate ./cmd/app/main.go

// 依赖下载
//go:generate go mod tidy

// gofmt goimports
//go:generate goimports -w .
//go:generate gofmt -l -w -s .

```


## 说明2
+ 目录说明
  + databse  //数据库连接工具
  + datas    //全局共享数据
  + errors   //自定义异常
  + generate //代码模板
  + schemas  //管理 generate 模板
  + tools    //代码中使用的部分工具
+ 扩展说明
  + 在 generate 中新增 .go
  + 在 init 函数中将自身注册到 schema 中
    + ```golang
      func init() {
         schemas.AddSchema("api", ApiSchema{})
      }
    ```
+ 后续(随缘后续)
  + 外挂 模板文件
  + 外挂 数据库字段影射 json 
  + 扩展支持数据库
