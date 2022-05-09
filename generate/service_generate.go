package generate

import (
	"fmt"
	"github.com/justxi1/bywin_generate/datas"
	"github.com/justxi1/bywin_generate/schemas"
	"github.com/justxi1/bywin_generate/tools"
	"os"
	"text/template"
)

const serviceTemplateStr = `
package service

import (
	"github.com/armory-toolkit/armory-web-starter-go/pkg/tools"
	"github.com/gin-gonic/gin"
	"{{.ProjectName}}/app/dao"
	"{{.ProjectName}}/app/model/entity"
	"{{.ProjectName}}/app/model/param"
	"{{.ProjectName}}/app/model/result"
)

var {{.TableUpperName}}ServiceInstance {{.TableUpperName}}Service

type {{.TableUpperName}}Service interface {
	Get{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}Search) (int64, []result.{{.TableUpperName}}ListSearch)
	Get{{.TableUpperName}}ById(ctx *gin.Context, id string) result.{{.TableUpperName}}Search
	Create{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}CreateForm) string
	Update{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}UpdateForm) string
	Delete{{.TableUpperName}}(ctx *gin.Context, id string) string
}

type {{.TableUpperName}}ServiceImpl struct {
	dao dao.{{.TableUpperName}}Dao
}

func New{{.TableUpperName}}ServiceImpl(dao dao.{{.TableUpperName}}Dao) {{.TableUpperName}}ServiceImpl {
	return {{.TableUpperName}}ServiceImpl{
		dao: dao,
	}
}

func (service {{.TableUpperName}}ServiceImpl) Get{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}Search) (int64, []result.{{.TableUpperName}}ListSearch) {
	count, searchResult := service.dao.Get{{.TableUpperName}}(ctx, param)
	r := make([]result.{{.TableUpperName}}ListSearch, 0)
	err := tools.CopySlices(&r, &searchResult)
	if err != nil {
		panic(err)
	}
	return count, r
}

func (service {{.TableUpperName}}ServiceImpl) Get{{.TableUpperName}}ById(ctx *gin.Context, id string) result.{{.TableUpperName}}Search {
	searchResult := service.dao.Get{{.TableUpperName}}ById(ctx, id)
	r := result.{{.TableUpperName}}Search{}
	err := tools.CopySlices(&r, &searchResult)
	if err != nil {
		panic(err)
	}
	return r
}

func (service {{.TableUpperName}}ServiceImpl) Create{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}CreateForm) string {
	{{.TableLowerName}} := entity.{{.TableUpperName}}{}
	err := tools.CopyProperties(&{{.TableLowerName}}, &param)
	if err != nil {
		panic(err)
	}
	return service.dao.Create{{.TableUpperName}}(ctx, {{.TableLowerName}})
}
func (service {{.TableUpperName}}ServiceImpl) Update{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}UpdateForm) string {
	{{.TableLowerName}} := entity.{{.TableUpperName}}{}
	err := tools.CopyProperties(&{{.TableLowerName}}, &param)
	if err != nil {
		panic(err)
	}
	return service.dao.Update{{.TableUpperName}}(ctx, {{.TableLowerName}})
}

func (service {{.TableUpperName}}ServiceImpl) Delete{{.TableUpperName}}(ctx *gin.Context, id string) string {
	return service.dao.Delete{{.TableUpperName}}(ctx, id)
}

func init() {
	{{.TableUpperName}}ServiceInstance = New{{.TableUpperName}}ServiceImpl(dao.{{.TableUpperName}}DaoInstance)
}

`

const serviceBasePath = "./app/service/"

func init() {
	schemas.AddSchema("service", ServiceSchema{})
}

type ServiceSchema struct {
}

func (s ServiceSchema) Start() {
	d := datas.Datas
	tools.CreateDirIfNotExists(serviceBasePath)
	for _, table := range d.TableInfos {
		filePath := serviceBasePath + table.TableName + "_service.go"
		if !tools.IsFileExists(filePath) {
			CreateService(filePath, table, d.ProjectName)
		} else if d.IsOverride {
			tools.DeleteFile(filePath)
			CreateService(filePath, table, d.ProjectName)
		}
	}
}

func CreateService(filePathStr string, table datas.TableInfo, ProjectName string) {
	tt := template.Must(template.New("service").Parse(serviceTemplateStr))
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
		"TableName":      table.TableName,
		"TableLowerName": table.TableLowerName,
		"TableUpperName": table.TableUpperName,
		"TableNameSnake": table.TableNameSnake,
		"TableComment":   table.TableComment,
		"ProjectName":    ProjectName,
	}
	_ = tt.Execute(file, vals)
	_ = file.Close()
}
