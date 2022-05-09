package schemas

var schemas map[string]Schema

func init() {
	schemas = make(map[string]Schema)
}

type Schema interface {
	Start()
}

func AddSchema(name string, s Schema) {
	schemas[name] = s
}

func StartSchemas() {
	for _, s := range schemas {
		s.Start()
	}
}
