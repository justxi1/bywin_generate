package generate

import (
	"fmt"
	"github.com/justxi1/bywingenerate/datas"
	"github.com/justxi1/bywingenerate/schemas"
	"github.com/justxi1/bywingenerate/tools"
	"os"
	"text/template"
)

const apiTemplateStr = `
package v1
import (
	"github.com/armory-toolkit/armory-web-starter-go/pkg/beans"
	"github.com/armory-toolkit/armory-web-starter-go/pkg/webstarter"
	"github.com/gin-gonic/gin"
	"{{.ProjectName}}/app/model/base_error"
	"{{.ProjectName}}/app/model/param"
	"{{.ProjectName}}/app/service"
	"net/http"
)

var {{.TableLowerName}}Service service.{{.TableUpperName}}Service

// Get{{.TableUpperName}}
// @Tags     {{.TableComment}}
// @Summary  查询 {{.TableComment}}
// @Accept   application/json
// @Param    param     query     param.{{.TableUpperName}}Search    true   "{{.TableComment}}查询接口"
// @Success  0 {object} beans.Result{data=result.{{.TableUpperName}}ListSearch}
// @Router   /v1/{{.TableNameSnake}} [get]
func Get{{.TableUpperName}}(ctx *gin.Context) {
	form := param.{{.TableUpperName}}Search{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		panic(baseError.ParamBindError)
	}
	count, result := {{.TableLowerName}}Service.Get{{.TableUpperName}}(ctx, form)
	ctx.JSON(http.StatusOK, beans.PageSuccess(count, result))
}

// Get{{.TableUpperName}}ById
// @Tags     {{.TableComment}}
// @Summary  查询 {{.TableComment}}
// @Accept   application/json
// @Param    id     path     string    true   "id"
// @Success  0 {object} beans.Result{data=result.{{.TableUpperName}}Search}
// @Router   /v1/{{.TableNameSnake}}/:id [get]
func Get{{.TableUpperName}}ById(ctx *gin.Context) {
	id := ctx.Param("id")
	result := {{.TableLowerName}}Service.Get{{.TableUpperName}}ById(ctx, id)
	ctx.JSON(http.StatusOK, beans.Success(result))
	
}

// Create{{.TableUpperName}}
// @Tags     {{.TableComment}}
// @Summary  创建组件模板
// @Accept   application/json
// @Param    param     body     param.{{.TableUpperName}}CreateForm    true   "{{.TableComment}}新增表单"
// @Success  0 {object} beans.Result{data=string}
// @Router   /v1/{{.TableNameSnake}} [post]
func Create{{.TableUpperName}}(ctx *gin.Context) {
	form := param.{{.TableUpperName}}CreateForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		panic(baseError.ParamBindError)
	}
	id := {{.TableLowerName}}Service.Create{{.TableUpperName}}(ctx, form)
	ctx.JSON(http.StatusOK, beans.Success(id))
}

// Update{{.TableUpperName}}
// @Tags     {{.TableComment}}
// @Summary  修改组件模板
// @Accept   application/json
// @Param    param     body     param.{{.TableUpperName}}UpdateForm     true   "组件修改表单"
// @Success  0 {object} beans.Result{data=string}
// @Router   /v1/{{.TableNameSnake}} [put]
func Update{{.TableUpperName}}(ctx *gin.Context) {
	form := param.{{.TableUpperName}}UpdateForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		panic(baseError.ParamBindError)
	}
	id := {{.TableLowerName}}Service.Update{{.TableUpperName}}(ctx, form)
	ctx.JSON(http.StatusOK, beans.Success(id))

}

// Delete{{.TableUpperName}}
// @Tags     {{.TableComment}}
// @Summary  删除{{.TableComment}}
// @Accept   application/json
// @Param    id       path     string     true   "id"
// @Success  0 {object} beans.Result{data=string}
// @Router   /v1/{{.TableNameSnake}}/:id [put]
func Delete{{.TableUpperName}}(ctx *gin.Context) {
	id := ctx.Param("id")
	result := {{.TableLowerName}}Service.Delete{{.TableUpperName}}(ctx, id)
	ctx.JSON(http.StatusOK, beans.Success(result))
}

func init() {
	{{.TableLowerName}}Service = service.{{.TableUpperName}}ServiceInstance
	g, err := webstarter.NewTimeOutGroup("/v1/{{.TableNameSnake}}")
	if err != nil {
		panic(err)
	}
	g.GET("/", Get{{.TableUpperName}})
	g.GET("/:id", Get{{.TableUpperName}}ById)
	g.POST("", Create{{.TableUpperName}})
	g.PUT("", Update{{.TableUpperName}})
	g.DELETE("/:id", Delete{{.TableUpperName}})
}
`
const apiBasePath = "./app/api/v1/"

func init() {
	schemas.AddSchema("api", ApiSchema{})
}

type ApiSchema struct {
}

func (s ApiSchema) Start() {
	d := datas.Datas
	tools.CreateDirIfNotExists(apiBasePath)
	for _, table := range d.TableInfos {
		filePath := apiBasePath + table.TableName + "_api.go"
		if !tools.IsFileExists(filePath) {
			CreateApi(filePath, table, d.ProjectName)
		} else if d.IsOverride {
			tools.DeleteFile(filePath)
			CreateApi(filePath, table, d.ProjectName)
		}
	}

}

func CreateApi(filePathStr string, table datas.TableInfo, ProjectName string) {
	tt := template.Must(template.New("api").Parse(apiTemplateStr))
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
