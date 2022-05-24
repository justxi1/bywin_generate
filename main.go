//go:generate go env -w GOPRIVATE=github.com/armory-toolkit
//go:generate go env -w GOINSECURE=github.com/armory-toolkit
//go:generate go get -u github.com/armory-toolkit/armory-web-starter-go@dev
//go:generate go mod tidy
//go:generate bywin_generate --type=mysql --override --projectName=generate_localTest --ip=xxx.xxx.xxx.138 --port=32556 --schema=test --user=root --password=root
//go:generate go generate ./cmd/app/main.go
//go:generate gofmt -w
package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql" //导入mysql包
	"github.com/justxi1/bywin_generate/database"
	"github.com/justxi1/bywin_generate/datas"
	_ "github.com/justxi1/bywin_generate/generate"
	"github.com/justxi1/bywin_generate/schemas"
)

var (
	dbType      string
	ip          string
	port        string
	schema      string
	user        string
	password    string
	isOverride  bool
	projectName string
)

func init() {
	flag.StringVar(&dbType, "dbType", "mysql", "database type （mysql）")
	flag.StringVar(&ip, "ip", "127.0.0.1", "mysql ip")
	flag.StringVar(&port, "port", "3306", "mysql port")
	flag.StringVar(&schema, "schema", "test", "database")
	flag.StringVar(&user, "user", "root", "user")
	flag.StringVar(&password, "password", "root", "password")
	flag.BoolVar(&isOverride, "override", false, "is override")
	flag.StringVar(&projectName, "projectName", "project_name", "is override")
}

func main() {
	flag.Parse()
	searcher := database.GetDatabaseSearch(dbType, ip, port, schema, user, password)
	tablse := searcher.SearchTables()

	d := datas.GlobalData{}
	d.IsOverride = isOverride
	d.TypeMapping = map[string]string{
		"varchar":    "string",
		"date":       "time.Time",
		"datetime":   "time.Time",
		"float":      "float32",
		"double":     "float64",
		"json":       "string",
		"bigint":     "int64",
		"mediumtext": "string",
	}
	d.ProjectName = projectName
	d.TableInfos = tablse
	datas.Datas = d
	schemas.StartSchemas()

}
