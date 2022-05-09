//go:generate go env -w GOPRIVATE=github.com/armory-toolkit
//go:generate go env -w GOINSECURE=github.com/armory-toolkit
//go:generate go get -u github.com/armory-toolkit/armory-web-starter-go@dev
//go:generate go mod tidy
//go:generate generate xxxxx
//go:generate go generate ./cmd/app/main.go
//go:generate gofmt -w
package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //导入mysql包
	"github.com/justxi1/bywin_generate/datas"
	_ "github.com/justxi1/bywin_generate/generate"
	"github.com/justxi1/bywin_generate/schemas"
	"github.com/justxi1/bywin_generate/tools"
)

var tabelSql = `
	select TABLE_NAME as "tableName", TABLE_COMMENT "tableComment"
	from TABLES
	WHERE TABLE_SCHEMA  = '%s' 

`
var columnSql = `
select 
	COLUMN_NAME as "columnName", 
	DATA_TYPE as "dataType", 
	COLUMN_COMMENT as "columnComment",
	COLUMN_KEY  as "columnKey"
from COLUMNS
where TABLE_NAME  = '%s' and TABLE_SCHEMA = '%s'

`
var (
	ip          string
	port        string
	schema      string
	user        string
	password    string
	isOverride  bool
	projectName string
)

const (
	baseUrl    string = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True"
	infoSchema string = "information_schema"
)

func init() {
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
	url := fmt.Sprintf(baseUrl, user, password, ip, port, infoSchema)
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic("数据库链接错误" + err.Error())
	}
	defer func() {
		if db != nil {
			_ = db.Close()
		}
	}()
	d := datas.GlobalData{}
	d.IsOverride = isOverride
	d.TypeMapping = map[string]string{
		"varchar": "string",
	}
	d.ProjectName = projectName
	d.TableInfos = searchTables(schema, db)
	datas.Datas = d
	schemas.StartSchemas()

}

func searchTables(schema string, db *sql.DB) []datas.TableInfo {
	r, err := db.Query(fmt.Sprintf(tabelSql, schema))
	if err != nil {
		panic(err)
	}
	defer func() {
		if r != nil {
			_ = r.Close()
		}
	}()
	result := make([]datas.TableInfo, 0)
	for r.Next() {
		tableName := ""
		tableComment := ""
		err := r.Scan(&tableName, &tableComment)
		if err != nil {
			panic(err)
		}
		info := datas.TableInfo{
			TableName:      tableName,
			TableLowerName: tools.LowerString(tableName),
			TableUpperName: tools.UpperString(tableName),
			TableNameSnake: tools.SnakeString(tableName),
			TableComment:   tableComment,
			Column:         searchColumns(schema, tableName, db),
		}
		result = append(result, info)
	}
	return result
}

func searchColumns(schema, tableName string, db *sql.DB) []datas.TableColumn {
	r, err := db.Query(fmt.Sprintf(columnSql, tableName, schema))
	if err != nil {
		panic(err)
	}
	defer func() {
		if r != nil {
			_ = r.Close()
		}
	}()

	result := make([]datas.TableColumn, 0)
	for r.Next() {
		columnName := ""
		dataType := ""
		columnComment := ""
		columnKey := ""
		err := r.Scan(&columnName, &dataType, &columnComment, &columnKey)
		if err != nil {
			panic(err)
		}
		column := datas.TableColumn{
			Name:      columnName,
			NameUpper: tools.UpperString(columnName),
			NameLower: tools.LowerString(columnName),
			NameSnake: tools.SnakeString(columnName),
			Type:      dataType,
			ColumKey:  columnKey,
			Comment:   columnComment,
		}
		result = append(result, column)
	}
	return result

}
