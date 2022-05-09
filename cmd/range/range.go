package main

import (
	"html/template"
	"os"
)

type EntetiesClass struct {
	Name  string
	Value int32
}

// In the template, we use rangeStruct to turn our struct values
// into a slice we can iterate over
const htmlTemplate = `
{{range $index, $element := .column}}
	{{.Name}}: {{.Type}}
{{end}}
`

type ColumnStruct struct {
	Name string
	Type string
}

func main() {
	//data := map[string][]EntetiesClass{
	//	"Yoga":    {{"Yoga1", 13}, {"Yoga2", 15}},
	//	"Pilates": {{"Pilates1", 3}, {"Pilates2", 6}, {"Pilates3", 9}},
	//}
	data := map[string]any{
		"column": []ColumnStruct{
			{Name: "Id", Type: "string"}, {Name: "Name", Type: "string"},
		},
	}

	t := template.New("t")
	t, err := t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
