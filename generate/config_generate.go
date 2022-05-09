package generate

import (
	"fmt"
	"github.com/justxi1/bywingenerate/schemas"
	"github.com/justxi1/bywingenerate/tools"
	"os"
	"text/template"
)

var yamlTemplateStr = `[log]
# 枚举| debug, info, warn, error, fatal, panic 默认 Info
level = "info"
#order = ["level", "time", "caller", "message"]
# 日志颜色
nocolor = false
# 输出日志路径
output = "./test.log"
span_name = "{{.ProjectName}}"

[web]
port = 8081
context_path = "/"
# 是否打印错误栈
#print_stack = false
# http 请求超时时间 默认 2000 毫秒
#timeout = 2000
# gin 模式| debug, release, test, 默认 release
#model = "debug"

[[gorm]]
host = "localhost"
port = 49157
database = "test"
username = "postgres"
password = "root"
max_idle_conn = 20
max_connect_num = 50
max_connect_life_time = 50
# 日志等级 info error warn Silent
log_level = "Info"
# postgres ssl_mode
ssl_mode = "disable"
database_type = "postgres"
# 查询超时 10s
sql_time_out = 10000
#慢查询 200 ms
slow_time = 200
#是否忽略 记录未找到异常
ignore_record_not_found = true

#[[gorm]]
#host = "localhost"
#port = 13306
#database = "test"
#user = "root"
#password = "root"
#max_idle_conn = 20
#max_connect_num = 50
#log_level = "Info"
#ssl_mode = "disable"
#database_type = "mysql"
#slow_time = 10


`

const configBasePath = "./config/"

func init() {
	schemas.AddSchema("config", ConfigSchema{})
}

type ConfigSchema struct {
}

func (s ConfigSchema) Start() {
	tools.CreateDirIfNotExists(configBasePath)
	filePathStr := configBasePath + "config.toml"
	tt := template.Must(template.New("config").Parse(yamlTemplateStr))
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
