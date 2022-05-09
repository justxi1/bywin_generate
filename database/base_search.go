package database

import (
	"github.com/justxi1/bywin_generate/datas"
	"strings"
)

type BaseSearch interface {
	SearchTables() []datas.TableInfo
}

const (
	typeMysql = "mysql"
)

func GetDatabaseSearch(dbType string, ip, port, schema, user, password string) BaseSearch {
	switch strings.ToLower(dbType) {
	case typeMysql:
		return NewMysqlSearch(ip, port, schema, user, password)
	default:
		return NewMysqlSearch(ip, port, schema, user, password)
	}
}
