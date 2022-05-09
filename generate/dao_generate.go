package generate

import (
	"fmt"
	"github.com/justxi1/bywin_generate/datas"
	"github.com/justxi1/bywin_generate/schemas"
	"github.com/justxi1/bywin_generate/tools"
	"os"
	"text/template"
)

const daoTemplateStr = `
package dao

import (
	"github.com/armory-toolkit/armory-web-starter-go/pkg/conn"
	"github.com/armory-toolkit/armory-web-starter-go/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"{{.ProjectName}}/app/model/base_error"
	"{{.ProjectName}}/app/model/entity"
	"{{.ProjectName}}/app/model/param"
)

var {{.TableUpperName}}DaoInstance {{.TableUpperName}}Dao

type {{.TableUpperName}}Dao interface {
	Get{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}Search) (int64, []entity.{{.TableUpperName}})
	Get{{.TableUpperName}}ById(ctx *gin.Context, id string) entity.{{.TableUpperName}}
	Create{{.TableUpperName}}(ctx *gin.Context, {{.TableLowerName}} entity.{{.TableUpperName}}) string
	Update{{.TableUpperName}}(ctx *gin.Context, {{.TableLowerName}} entity.{{.TableUpperName}}) string
	Delete{{.TableUpperName}}(ctx *gin.Context, id string) string
}

type {{.TableUpperName}}DaoImpl struct {
}

func New{{.TableUpperName}}Dao() {{.TableUpperName}}Dao {
	return {{.TableUpperName}}DaoImpl{}
}

func init() {
	{{.TableUpperName}}DaoInstance = New{{.TableUpperName}}Dao()
}

func (dao {{.TableUpperName}}DaoImpl) Get{{.TableUpperName}}(ctx *gin.Context, param param.{{.TableUpperName}}Search) (int64, []entity.{{.TableUpperName}}) {
	var count int64
	result := make([]entity.{{.TableUpperName}}, 0)
	if tx := conn.Connect().
		Table(entity.{{.TableUpperName}}{}.TableName()).
		Count(&count); tx.Error != nil {
		panic(baseError.DatabaseOptError.DetailError(tx.Error.Error()))
	}
	if count == 0 {
		return count, result
	}

	page, limit := tools.PageInfo(ctx)
	if tx := conn.Connect().
		Table(entity.{{.TableUpperName}}{}.TableName()).
		Offset((page - 1) * limit).Limit(limit).
		Find(&result); tx.Error != nil {
		panic(baseError.DatabaseOptError.DetailError(tx.Error.Error()))
	}
	return count, result
}

func (dao {{.TableUpperName}}DaoImpl) Get{{.TableUpperName}}ById(ctx *gin.Context, id string) entity.{{.TableUpperName}} {
	result := entity.{{.TableUpperName}}{}
	if tx := conn.Connect().
		Table(entity.{{.TableUpperName}}{}.TableName()).
		Where("id", id).
		Find(&result); tx.Error != nil {
		panic(baseError.DatabaseOptError.DetailError(tx.Error.Error()))
	}
	if result.Id == "" || result.Id != id {
		panic(baseError.DatabaseOptError.DetailError("数据未找到"))
	}
	return result
}

func (dao {{.TableUpperName}}DaoImpl) Create{{.TableUpperName}}(ctx *gin.Context, {{.TableLowerName}} entity.{{.TableUpperName}}) string {
	{{.TableLowerName}}.Id = uuid.NewString()
	if tx := conn.Connect().
		Table(entity.{{.TableUpperName}}{}.TableName()).
		Create(&{{.TableLowerName}}); tx.Error != nil {
		panic(baseError.DatabaseOptError.DetailError(tx.Error.Error()))
	}
	return {{.TableLowerName}}.Id
}

func (dao {{.TableUpperName}}DaoImpl) Update{{.TableUpperName}}(ctx *gin.Context, {{.TableLowerName}} entity.{{.TableUpperName}}) string {
	if tx := conn.Connect().
		Table(entity.{{.TableUpperName}}{}.TableName()).
		Updates(&{{.TableLowerName}}); tx.Error != nil {
		panic(baseError.DatabaseOptError.DetailError(tx.Error.Error()))
	}
	return {{.TableLowerName}}.Id
}

func (dao {{.TableUpperName}}DaoImpl) Delete{{.TableUpperName}}(ctx *gin.Context, id string) string {
	if tx := conn.Connect().Table(entity.{{.TableUpperName}}{}.TableName()).
		Where("id", id).
		Delete(entity.{{.TableUpperName}}{}); tx.Error != nil {
		panic(tx.Error)
	}
	return id
}
`

const daoBasePath = "./app/dao/"

func init() {
	schemas.AddSchema("dao", DaoSchema{})
}

type DaoSchema struct {
}

func (s DaoSchema) Start() {
	d := datas.Datas
	tools.CreateDirIfNotExists(daoBasePath)
	for _, table := range d.TableInfos {
		filePath := daoBasePath + table.TableName + "_dao.go"
		if !tools.IsFileExists(filePath) {
			CreateDao(filePath, table, d.ProjectName)
		} else if d.IsOverride {
			tools.DeleteFile(filePath)
			CreateDao(filePath, table, d.ProjectName)
		}
	}
}

func CreateDao(filePathStr string, table datas.TableInfo, ProjectName string) {
	tt := template.Must(template.New("dao").Parse(daoTemplateStr))
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
