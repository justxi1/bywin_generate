package datas

var Datas GlobalData

type GlobalData struct {
	IsOverride  bool
	ProjectName string
	TypeMapping map[string]string
	TableInfos  []TableInfo
}

type TableInfo struct {
	TableName      string
	TableUpperName string
	TableLowerName string
	TableNameSnake string
	TableComment   string
	Column         []TableColumn
}

type TableColumn struct {
	Name      string
	NameUpper string
	NameLower string
	NameSnake string
	Type      string
	ColumKey  string
	Comment   string
}
