// 从第二个参数开始就是需要生成的队列的值类似，比如 ./queue int string MyInt
package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

var tpl = `// Code generated, DO NOT EDIT.
package {{.Package}}

type {{.MyType}}Queue struct {
    q []{{.MyType}}
}

func New{{.MyType}}Queue() *{{.MyType}}Queue {
    return &{{.MyType}}Queue{
        q: []{{.MyType}}{},
    }
}

func (o *{{.MyType}}Queue) Insert(v {{.MyType}}) {
    o.q = append(o.q, v)
}

func (o *{{.MyType}}Queue) Remove() {{.MyType}} {
    if len(o.q) == 0 {
        panic("empty queue")
    }
    first := o.q[0]
    o.q = o.q[1:]
    return first
}
`

func main() {
	tt := template.Must(template.New("api").Parse(tpl))
	for i := 1; i < len(os.Args); i++ { // 读入传参
		var (
			dest = strings.ToLower(os.Args[i]) + "_queue.go"
			err  error
			file *os.File
		)
		if file, err = os.Create(dest); err != nil {
			if !os.IsExist(err) {
				fmt.Printf("Could not create %s: %s (skip)\n", dest, err)
				continue
			}
			_ = os.Remove(dest)
		}
		vals := map[string]string{
			"MyType":  os.Args[i],
			"Package": os.Getenv("GOPACKAGE"),
		}
		_ = tt.Execute(file, vals)
		_ = file.Close()
	}
}
