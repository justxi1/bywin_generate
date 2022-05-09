package generate

import (
	"fmt"
	"github.com/justxi1/bywingenerate/datas"
	"github.com/justxi1/bywingenerate/schemas"
	"github.com/justxi1/bywingenerate/tools"
	"os"
	"text/template"
)

const mainTemplateStr = `
//go:generate swag init  --parseDependency --parseDepth 3 -g ./main.go -o ../../doc/swagger
package main

import (
	"github.com/armory-toolkit/armory-web-starter-go/pkg/webstarter"
	_ "{{.ProjectName}}/app/api/v1"
	_ "{{.ProjectName}}/doc/swagger"
)

// @title {{.ProjectName}}
// @version 1.0
// @description {{.ProjectName}} 接口
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// 指令 swag init  --parseDependency --parseDepth 3 -g ./cmd/app/Main.go -o ./doc/swagger
// 打包指令:
// go env -w GOOS=linux
// go build -o {{.ProjectName}} ./cmd/app/Main.go
// 启动环境变量配置 CONFIG_FILE = ./config/config.toml
func main() {
	webstarter.StartWebServer()
}	
`

const mainBasePath = "./cmd/app/"

func init() {
	schemas.AddSchema("main", MainSchema{})
}

type MainSchema struct {
}

func (s MainSchema) Start() {
	d := datas.Datas
	tools.CreateDirIfNotExists(mainBasePath)
	filePathStr := mainBasePath + "main.go"
	tt := template.Must(template.New("main").Parse(mainTemplateStr))
	var err error
	var file *os.File
	if file, err = os.Create(filePathStr); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Could not create %s: %s (skip)\n", filePathStr, err)
			return
		}
		_ = os.Remove(filePathStr)
	}
	vals := map[string]string{
		"ProjectName": d.ProjectName,
	}
	_ = tt.Execute(file, vals)
	_ = file.Close()

}
