package generate

import (
	"fmt"
	"github.com/justxi1/bywin_generate/datas"
	"github.com/justxi1/bywin_generate/schemas"
	"github.com/justxi1/bywin_generate/tools"
	"os"
	"text/template"
)

const modelTemplateStr string = `
package entity

import (
	"time"
)
type {{.TableUpperName}} struct {
	{{range $i, $column := .Columns}}
	{{.Name}} {{.Type}} {{.Tag}} // {{.Comment}}{{end}}
}

func ({{.TableUpperName}}) TableName() string {
	return "{{.TableName}}"
}
`
const paramTemplateStr string = `
package param

//{{.TableUpperName}}Search 查询
type {{.TableUpperName}}Search struct {
	{{range $i, $column := .Columns}}
	{{.Name}} {{.Type}} {{.FormTag}} // {{.Comment}}{{end}}
}

//{{.TableUpperName}}CreateForm  创建
type {{.TableUpperName}}CreateForm struct {
	{{range $i, $column := .Columns}}
	{{.Name}} {{.Type}} {{.FormTag}} // {{.Comment}}{{end}}

}

//{{.TableUpperName}}UpdateForm  修改
type {{.TableUpperName}}UpdateForm struct {
	{{range $i, $column := .Columns}}
	{{.Name}} {{.Type}} {{.FormTag}} // {{.Comment}}{{end}}
}
`
const resultTemplateStr string = `
package result

type {{.TableUpperName}}ListSearch struct {
	{{range $i, $column := .Columns}}
	{{.Name}} {{.Type}} {{.JsonTag}} // {{.Comment}}{{end}}
}

type {{.TableUpperName}}Search struct {
	{{range $i, $column := .Columns}}
	{{.Name}} {{.Type}} {{.JsonTag}} // {{.Comment}}{{end}}
}

`

const (
	pkiTag         = "`json:\"%s\" gorm:\"primaryKey;column:%s\"`"
	baseTag        = "`json:\"%s\" gorm:\"column:%s\"`"
	jsonTag        = "`json:\"%s\"`"
	formTag        = "`json:\"%s\" form:\"%s\"`"
	PRI            = "PRI"
	modelBasePath  = "./app/model/entity/"
	paramBasePath  = "./app/model/param/"
	resultBasePath = "./app/model/result/"
)

func init() {
	schemas.AddSchema("model", ModelsSchema{})
}

type ModelsSchema struct {
}

func (ModelsSchema) Start() {
	d := datas.Datas
	tools.CreateDirIfNotExists(modelBasePath)
	tools.CreateDirIfNotExists(paramBasePath)
	tools.CreateDirIfNotExists(resultBasePath)
	for _, tableInfo := range d.TableInfos {
		columnStructs := make([]ModelColumnStructs, 0)
		for _, column := range tableInfo.Column {
			columnStructs = append(columnStructs, NewModelColumnStructs(column, d.TypeMapping))
			createModel(columnStructs, modelBasePath+tableInfo.TableName+".go", tableInfo, "model", modelTemplateStr, d.IsOverride)
			createModel(columnStructs, paramBasePath+tableInfo.TableName+".go", tableInfo, "param", paramTemplateStr, d.IsOverride)
			createModel(columnStructs, resultBasePath+tableInfo.TableName+".go", tableInfo, "result", resultTemplateStr, d.IsOverride)
		}
	}
}
func NewModelColumnStructs(column datas.TableColumn, mapping map[string]string) ModelColumnStructs {
	return ModelColumnStructs{
		Name:    column.NameUpper,
		Type:    mapping[column.Type],
		Tag:     getTags(column),
		JsonTag: getJsonTags(column),
		FormTag: getFormTags(column),
		Comment: column.Comment,
	}
}

func createModel(columns []ModelColumnStructs, filePathStr string, table datas.TableInfo, templateName, templateStr string, IsOverride bool) {
	if tools.IsFileExists(filePathStr) && !IsOverride {
		return
	}

	tt := template.Must(template.New(templateName).Parse(templateStr))
	var err error
	var file *os.File
	if file, err = os.Create(filePathStr); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Could not create %s: %s (skip)\n", filePathStr, err)
			return
		}
		_ = os.Remove(filePathStr)
	}
	vals := map[string]any{
		"TableUpperName": table.TableUpperName,
		"Columns":        columns,
	}
	_ = tt.Execute(file, vals)
	_ = file.Close()
}

func getTags(column datas.TableColumn) string {
	if column.ColumKey == PRI {
		return fmt.Sprintf(pkiTag, column.NameSnake, column.Name)
	}
	return fmt.Sprintf(baseTag, column.NameSnake, column.Name)
}

func getJsonTags(column datas.TableColumn) string {
	return fmt.Sprintf(jsonTag, column.NameSnake)
}
func getFormTags(column datas.TableColumn) string {
	return fmt.Sprintf(formTag, column.NameSnake, column.NameSnake)
}

type ModelColumnStructs struct {
	Name    string
	Type    string
	Tag     string
	JsonTag string
	FormTag string
	Comment string
}
