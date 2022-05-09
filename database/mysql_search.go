package database

import (
	"database/sql"
	"fmt"
	"github.com/justxi1/bywin_generate/datas"
	"github.com/justxi1/bywin_generate/tools"
)

const tabelSql string = `
	select TABLE_NAME as "tableName", TABLE_COMMENT "tableComment"
	from TABLES
	WHERE TABLE_SCHEMA  = '%s' 

`
const columnSql string = `
select 
	COLUMN_NAME as "columnName", 
	DATA_TYPE as "dataType", 
	COLUMN_COMMENT as "columnComment",
	COLUMN_KEY  as "columnKey"
from COLUMNS
where TABLE_NAME  = '%s' and TABLE_SCHEMA = '%s'

`

const (
	baseUrl    string = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True"
	infoSchema string = "information_schema"
)

type MysqlSearch struct {
	db     *sql.DB
	schema string
}

func NewMysqlSearch(ip, port, schema, user, password string) MysqlSearch {
	url := fmt.Sprintf(baseUrl, user, password, ip, port, infoSchema)
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic("数据库链接错误" + err.Error())
	}
	return MysqlSearch{
		db:     db,
		schema: schema,
	}
}

func (d MysqlSearch) SearchTables() []datas.TableInfo {
	r, err := d.db.Query(fmt.Sprintf(tabelSql, d.schema))
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
			Column:         searchColumns(d.schema, tableName, d.db),
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
