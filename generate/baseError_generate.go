package generate

import (
	"fmt"
	"github.com/justxi1/bywin_generate/schemas"
	"github.com/justxi1/bywin_generate/tools"
	"os"
	"text/template"
)

const baseErrorTemplate = `
package baseError

import (
	"github.com/armory-toolkit/armory-web-starter-go/pkg/beans"
	"net/http"
)

var ParamBindError = beans.BaseError{HttpCode: http.StatusInternalServerError, Code: 1, Msg: "参数绑定异常"}
var DatabaseOptError = beans.BaseError{HttpCode: http.StatusInternalServerError, Code: 1, Msg: "数据库操作异常"}

`

const baseErrorBasePath = "./app/model/base_error/"

func init() {
	schemas.AddSchema("baseError", BaseErrorSchema{})
}

type BaseErrorSchema struct {
}

func (s BaseErrorSchema) Start() {
	tools.CreateDirIfNotExists(baseErrorBasePath)
	filePathStr := baseErrorBasePath + "base_error.go"
	tt := template.Must(template.New("baseError").Parse(baseErrorTemplate))
	var err error
	var file *os.File
	if file, err = os.Create(filePathStr); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Could not create %s: %s (skip)\n", filePathStr, err)
			return
		}
		_ = os.Remove(filePathStr)
	}
	vals := map[string]string{}
	_ = tt.Execute(file, vals)
	_ = file.Close()

}
